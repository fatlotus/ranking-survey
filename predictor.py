import numpy
import math
import random

class Predictor(object):
    def __init__(self, users, things):
        """
        Initialize this predictor given a list of users and things.
        """

        self.users = list(users)
        self.things = list(things)

        shape = (len(self.users), len(self.things))
        self.x = numpy.asmatrix(numpy.zeros(shape))
        self.xp = numpy.asmatrix(numpy.zeros(shape))
        self.z = numpy.asmatrix(numpy.zeros(shape))

        self.nu = 1
        self.alpha = 1
        self.lambd = 0.04
        self.t = 1

        self.samples = []

    def _hinge_loss(self, samples):
        """
        Computes the hinge loss for the given samples.
        """

        sum = 0.0
        for (user, a, b) in samples:
            sum += max(1 - (self.x[user, a] - self.x[user, b]), 0)
        return sum / float(len(samples))

    def _gradient_loss(self, samples):
        """
        Jankily computes the gradient of the hinge loss of the belief matrix 
        on the given samples.
        """

        result = numpy.zeros((len(self.users), len(self.things)))
        before = self._hinge_loss(samples)
        eps = 0.0001

        for i, _ in enumerate(self.users):
            for j, _ in enumerate(self.things):
                self.x[i, j] += eps
                result[i, j] = (self._hinge_loss(samples) - before) / eps
                self.x[i, j] -= eps

        return result

    def _update(self, samples):
        """
        Updates the given belief matrix after having recevied a new sample.
        """

        alphaP = (1 + math.sqrt(1 + 4*self.alpha**2)) / 2
        u, s, v = numpy.linalg.svd(self.z -
            self._gradient_loss(samples) * self.nu)
        for i in xrange(len(s)):
            s[i] = max(0, s[i] - self.lambd)

        self.xp = self.x
        self.x = u * numpy.diag(s) * v
        self.z = self.x + (self.x - self.xp) * (self.alpha - 1) / alphaP
        self.alpha = alphaP

    def generate(self, mask=lambda user, a, b: (a != b)):
        """
        Generates a tuple (user, a, b) for the ideal next option.

        Callers can optionally specify a boolean function that decides whether
        a given pairing is allowed to be considered. This could, for instance,
        select just those options that are valid for the current user.
        """

        candidates = []
        total = 0.0
        for u, user in enumerate(self.users):
            for a, thinga in enumerate(self.things):
                for b, thingb in enumerate(self.things):
                    exponent = abs(self.x[u, a] - self.x[u, b])
                    weight = math.exp(-exponent / self.t)
                    if not mask(user, thinga, thingb):
                        continue
                    candidates.append((weight, user, thinga, thingb))
                    total += weight

        value = random.uniform(0, total)
        for (weight, user, thinga, thingb) in candidates:
            value -= weight
            if value <= 0:
                return (user, thinga, thingb)

    def respond(self, user, thinga, thingb):
        """
        Informs this predictor that, of the two options, the user prefers the
        first one (A) more than (B).
        """

        self.samples.append((
            self.users.index(user),
            self.things.index(thinga),
            self.things.index(thingb)
        ))
        self._update(self.samples)

    def ranking(self, user):
        """
        Returns the predicted ranking for a given user, in the form
        (option, score) where score is a rough magnitude for the user's 
        confidence in this position.
        """

        u, ary = self.users.index(user), numpy.asarray(self.x)
        indices = sorted(enumerate(ary[u,:]), key=lambda (i, x): -x)
        return [(self.things[i], x) for i, x in indices]

def main():
    predictor = Predictor(
        ["Subject 1", "Subject 2", "Subject 3", "Subject 4"],
        ["Apples", "Bananas", "Cantelope", "Date"]
    )
    simulation = lambda user, a, b: (min(a, b), max(a, b))

    for i in xrange(100):
        user, a, b = predictor.generate()
        print("Ask {} to compare {} or {}".format(user, a, b))

        a, b = simulation(user, a, b)
        print("They prefer {} over {}".format(a, b))

        predictor.respond(user, a, b)
        print("Now we have: {!r}".format(
            [x[0] for x in predictor.ranking(user)]))

if __name__ == "__main__":
    main()
