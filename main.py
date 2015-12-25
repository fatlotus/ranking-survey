import sys, os; sys.path.append(os.path.dirname(__file__) + "/vendor")

from google.appengine.ext import ndb

from flask import Flask, session, redirect, render_template, url_for, request

import random
from markupsafe import Markup
import json

from predictor import Predictor

app = Flask(__name__)

def generate_color():
    h = random.uniform(0, 360)
    s = random.uniform(0.2, 1.0)
    l = random.uniform(0.2, 0.8)
    return (h, s, l)

@app.template_global()
def show_color((h, s, l)):
    return Markup("""<div class='box' style='background-color:
                     hsl({:}, {:%}, {:%})'></div>""".format(h, s, l))

@app.template_global()
def dumps(value):
    return json.dumps(value)

class Survey(ndb.Model):
    predictor = ndb.PickleProperty()

    @classmethod
    def generate(klass, key, size):
        users = [(generate_color(), generate_color()) for i in xrange(size)]
        things = [generate_color() for i in xrange(size)]

        return Survey.get_or_insert(key, predictor=Predictor(users, things))

# Allow Surveys in URL parameters.
from werkzeug.routing import BaseConverter
class SurveyConverter(BaseConverter):
    def to_python(self, value):
        return Survey.generate(value, 5)

    def to_url(self, value):
        return value.key.id()

app.url_map.converters["survey"] = SurveyConverter

@app.route("/")
def home():
    return redirect("/" + str(random.randint(0, 10 ** 10)))

@app.route("/")
def home_page():
    return redirect("/" + str(random.randint(0, 10**10)))

# Allow users to view each survey.
@app.route("/<survey:survey>")
def display_question(survey):
    pred = survey.predictor
    user, a, b = pred.generate()
    rankings = [(u, pred.ranking(u)) for u in pred.users]
    return render_template("prompt.html",
        user=user, a=a, b=b, rankings=rankings)

# Allow users to run surveys.
def tupalov(value):
    """
    Converts the given JSON object back to an immutable type.
    """

    if type(value) is list:
        return tuple([tupalov(x) for x in value])
    else:
        return value

@app.route("/<survey:survey>", methods=["POST"])
def prompt(survey):
    a = tupalov(json.loads(request.form["a"]))
    b = tupalov(json.loads(request.form["b"]))
    user = tupalov(json.loads(request.form["user"]))
    choice = tupalov(json.loads(request.form["choice"]))

    if choice == b:
        a, b = b, a

    print user, a, b

    survey.predictor.respond(user, a, b)
    survey.put()

    return redirect(url_for("display_question", survey=survey))
