import json
import random

videos = [
    '<strong>Video:</strong> Cat in Pirate Costume<br/><iframe src="https://www.youtube.com/embed/SCfcn1Rtqz0" frameborder="0" allowfullscreen></iframe>',
    '<strong>Video:</strong> Cat in Monkey Suit<br/><iframe src="https://www.youtube.com/embed/YtEhX7mhDJc?t=7s" frameborder="0" allowfullscreen></iframe>'
]

for i in xrange(100):
    question = {
        "survey": "survey",
        "choices": videos,
        "precision": random.choice([3, 8, 12]),
        "exclusive": random.choice([False, True]),
    }
    print json.dumps(question)
