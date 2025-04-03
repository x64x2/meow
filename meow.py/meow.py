import argparse
import cv2
import numpy as np
import numpy.typing as npt
import struct
from lxml import etree

import tqdm

import matplotlib.pyplot as plt

class Rectangle:
    def __init__(
        self,
        x0: int,
        y0: int,
        width: int,
        height: int,
        weight: int
    ):
        self.x0 = x0
        self.y0 = y0
        self.width = width
        self.height = height
        self.weight = weight

    def __call__(
        self,
        image: npt.NDArray[np.int_],
        window: tuple[int, int, float]
    ) -> float:
        x0 = int(self.x0 * window[2] + window[0])
        y0 = int(self.y0 * window[2] + window[1])
        x1 = x0 + int(self.width * window[2])
        y1 = y0 + int(self.height * window[2])

        value = image[y1, x1] - image[y0, x1] - image[y1, x0] + image[y0, x0]

        return self.weight * value

    def fromXML(rectElement: etree._Element):
        line = rectElement.text

        if line.endswith('.'):
            line = line[:-1]

        x0, y0, width, height, weight = [int(_) for _ in line.split(' ')[-5:]]

        return Rectangle(x0, y0, width, height, weight)

class WeakClassifier:
    """
    A Haar-like feature classifier
    """

    def __init__(
        self,
        feature: list[Rectangle],
        activations: tuple[float, float],
        threshold: float,
    ):
        self.feature = feature
        self.activations = activations
        self.threshold = threshold #################################

    def __call__(
        self,
        image: npt.NDArray[np.int_],
        window: tuple[int, int, float]
    ) -> float:
        value = sum(map(lambda _: _(image, window), self.feature))

        return self.activations[1] if value <= self.threshold else self.activations[0]
    def fromXML(weakClassifierElement: etree._Element, featureElement: etree._Element):
        feature = [
            Rectangle.fromXML(rectElement)
            for rectElement
            in featureElement.findall('./rects/_')
        ]

        threshold = float(
            weakClassifierElement.find('./internalNodes').text.split(' ')[-1]
        )

        activations = tuple(
            [
                float(_)
                for _
                in weakClassifierElement.find('./leafValues').text.split(' ')[-2:]
            ]
        )

        return WeakClassifier(feature, activations, threshold)

class StrongClassifier:
    """
    A set of weak features, for together they are powerful!
    """

    def __init__(self, classifiers: list[WeakClassifier], threshold: float):
        self.threshold = threshold * 0.6 ######################################
        self.classifiers = classifiers

    def __call__(
        self,
        image: npt.NDArray[np.int_],
        window: tuple[int, int, float]
    ) -> bool:
        classifyEach = lambda c: c(image, window)
        classifications = map(classifyEach, self.classifiers)
        return sum(classifications) <= self.threshold

    def fromXML(stageElement: etree._Element, featureElements: etree._Element):
        threshold = float(stageElement.find('./stageThreshold').text)

        weakClassifierElements = stageElement.findall('./weakClassifiers/_')

        classifiers = [
            WeakClassifier.fromXML(weakClassifierElement, featureElement)
            for weakClassifierElement, featureElement
            in zip(weakClassifierElements, featureElements)
        ]

        return StrongClassifier(classifiers, threshold)

class CascadeClassifier:
    """
    A series of strong classifiers, an image must pass all.
    """

    def __init__(
        self,
        stages: list[StrongClassifier],
        filterSize: int
    ):
        """
        Parameters
        ----------
        stages : List[StrongClassifier]
            The cascade stages
        filterSize : int
            The width and the height of the unscaled filter in pixels

        Returns
        -------
        None.
        """

        self.filterSize = filterSize
        self.stages = stages

    def __call__(
            self,
            image: npt.NDArray[np.uint8],
            window: tuple[int, int, float]
            ) -> bool:
        for i, stage in enumerate(self.stages):
            if not stage(image, window):
                #if i > 12: print("Blocked at: {}".format(i))
                return False

        return True

    def fromXML(cascadeElement: etree._Element):
        assert cascadeElement.find('./stageType').text == 'BOOST'
        assert cascadeElement.find('./featureType').text == 'HAAR'

        filterSize = int(cascadeElement.find('./width').text)
        assert filterSize == int(cascadeElement.find('./height').text)

        stageElements = cascadeElement.findall('./stages/_')
        assert len(stageElements) == int(cascadeElement.find('stageNum').text)

        featureElements = cascadeElement.findall('./features/_')

        stages = [
            StrongClassifier.fromXML(stageElement, featureElements)
            for stageElement
            in stageElements
        ]

        return CascadeClassifier(stages, filterSize)

if __name__ == '__main__':
    # Parse arguments.
    argsParser = argparse.ArgumentParser(
        description='Detect cat faces in images',
    )
    argsParser.add_argument(
        'classifier',
        type=argparse.FileType('rb'),
        help='An OpenCV cascade classifier XML file'
    )
    argsParser.add_argument(
        'image',
        type=argparse.FileType('rb'),
        help='The image to find the faces in'
    )
    args = argsParser.parse_args(
        [
            '/media/trisquel/pictures/meow/meow.py/haarcascade_frontalface_default.xml',
            '/media/trisquel/pictures/meow/meow.py/53_2.png'
        ]
    )

    # Parse the XML file.
    elementTree = etree.parse(args.classifier)
    cascadeElement = elementTree.find('./cascade')

    cascade = CascadeClassifier.fromXML(cascadeElement)

    # Open the image.
    image = np.frombuffer(args.image.read(), dtype=np.uint8)
    image = cv2.imdecode(image, cv2.IMREAD_GRAYSCALE)
    toplot = np.array(image)

    # Pre-process it.
    image = cv2.integral(image)

    image_dim = min(image.shape)
    for size in range(int(0.7 * image_dim), image_dim, int(0.05 * image_dim)):
        scale = size / cascade.filterSize
        for y in tqdm.trange(0, image.shape[0] - size, int(0.03 * image_dim)):
            for x in range(0, image.shape[1] - size, int(0.03 * image_dim)):
                if cascade(image, (x, y, scale)):
                    print(x, y, size)
                    toplot = cv2.rectangle(toplot, (x, y), (x + size, y + size), (255, 0, 0), 2)

    plt.imshow(toplot)
