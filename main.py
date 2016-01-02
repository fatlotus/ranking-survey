import sys
import os
sys.path.append(os.path.dirname(__file__) + "/vendor")

from google.appengine.ext import ndb

from flask import Flask, session, redirect, render_template, url_for, request

import json
import random
import math
import urlparse
from markupsafe import Markup

from predictor import Predictor

app = Flask(__name__)

YOUTUBE_VIDEOS = [
    ("PHAc3_MEjgQ", "Funny Cats Acting Like..."),
    ("YtEhX7mhDJc", "Cat + Monkey"),

    ("6ztm7YkLElI", "World's Most Amazing Dogs"),
    ("kMhw5MFYU0s", "Dogs Who Fail At"),

    ("2mUBHsxoK7I", "The Funniest Parrots ever"),
    ("Rpu-cNjT2Dg", "Funny Parrot Videos"),
]


def generate_color():
    """
    Generates a random HSL color.
    """

    h = random.uniform(0, 360)
    s = random.uniform(0.4, 1.0)
    l = random.uniform(0.4, 0.6)
    return (h, s, l)


def distance(a, b):
    """
    Computes the Euclidean distance between two HSL colors.
    """

    dh = (a[0] - b[0] + 180 + 360) % 360 - 180
    ds = a[1] - b[1]
    dl = a[2] - b[2]

    return math.sqrt(dh * dh + ds * ds + dl * dl)


def generate_pair():
    """
    Generates a pair of colors that are somewhat different.
    """

    a, b = generate_color(), generate_color()
    while distance(a, b) < 2.:
        a, b = generate_color(), generate_color()
    return a, b


@app.template_global()
def show_color(color):
    """
    Renders a color as an HSL box.
    """

    return Markup("""<div class='box' style='background-color:
                     hsl({:}, {:%}, {:%})'></div>""".format(*color))


@app.template_global()
def show_video(video):
    """
    Displays a video in an embed box.
    """

    return Markup(
        """<iframe width="560" height="315"
             src="https://www.youtube.com/embed/{0}">
           </iframe>""").format(*video)


@app.template_global()
def dumps(value):
    return json.dumps(value)


class Survey(ndb.Model):
    predictor = ndb.PickleProperty()
    suite = ndb.StringProperty()

    @classmethod
    def generate(klass, key, size, suite="colors"):
        if suite == "colors":
            users = [generate_pair() for i in xrange(size)]
            things = [generate_color() for i in xrange(size)]
        else:
            users = list(range(size))[:6]
            things = YOUTUBE_VIDEOS[:size]

        return Survey.get_or_insert(key, predictor=Predictor(users, things),
                                    suite=suite)

# Allow Surveys in URL parameters.
from werkzeug.routing import BaseConverter


class SurveyConverter(BaseConverter):

    """
    Converts URLs into Survey keys.

    e.x.
    /foobar               -> (3 people, 3 colors)
    /size=7               -> (7 people, 7 colors)
    /suite=youtube        -> (3 people, 3 videos)
    /suite=youtube&size=5 -> (5 people, 5 videos)
    """

    def to_python(self, value):
        params = urlparse.parse_qs(value)
        size = int(params.get("size", ["3"])[0])
        suite = params.get("suite", ["colors"])[0]
        if suite not in ("colors", "youtube"):
            suite = "colors"
        return Survey.generate(value, size, suite=suite)

    def to_url(self, value):
        return value.key.id()

app.url_map.converters["survey"] = SurveyConverter


@app.route("/")
def home():
    """
    Randomly generate a survey as demo.
    """

    return redirect("/" + str(random.randint(0, 10 ** 10)) + "&size=10")


@app.route("/<survey:survey>")
def display_question(survey):
    """
    Display a given survey, allowing users to take it.
    """

    pred = survey.predictor
    user = pred.users[int(request.args.get("user", "0"))]
    user, a, b = pred.generate(mask=lambda u, a, b: u == user)
    rankings = [(u, pred.ranking(u)) for u in pred.users]
    return render_template("prompt.html", suite=survey.suite,
                           user=user, a=a, b=b, rankings=rankings)


def tupalov(value):
    """
    Converts the given JSON object back to an immutable type.
    """

    if isinstance(value, list):
        return tuple([tupalov(x) for x in value])
    else:
        return value


@app.route("/<survey:survey>", methods=["POST"])
def prompt(survey):
    """
    Process a choice from a user.
    """

    a = tupalov(json.loads(request.form["a"]))
    b = tupalov(json.loads(request.form["b"]))
    user = tupalov(json.loads(request.form["user"]))
    choice = tupalov(json.loads(request.form["choice"]))

    if choice == b:
        a, b = b, a

    survey.predictor.respond(user, a, b)
    survey.put()

    return redirect(url_for("display_question", survey=survey))
