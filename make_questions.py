#!/usr/bin/env python


"""
Generates a series of questions in the surveys/chunk-NN.json folder.
"""


import collections
import json
import random
import itertools
import os
import shutil

Question = collections.namedtuple(
    "Question", ["survey", "choices", "precision", "exclusive"])


def generate_comparison(choices, survey):
    """
    Generate a comparison question for a particular survey.

    Args:
        choices: things to rate, as list of HTML strings
        survey: survey URL
    Returns:
        a Question
    """

    return Question(
        survey=survey,
        choices=tuple(sorted(random.sample(choices, 2))),
        precision=2,
        exclusive=True
    )


def generate_rating(choices, survey, precision):
    """
    Generate a comparison question for a particular survey.

    Args:
        choices: things to rate, as list of HTML strings
        survey: survey URL
        precision: max # of stars
    Returns:
        a Question
    """

    return Question(
        survey=survey,
        choices=(random.choice(choices),),
        precision=precision,
        exclusive=False,
    )


def repeat(count, func, *vargs, **kwargs):
    """
    Generate exactly count items from a distribution.

    Args:
        count: # of things to generate
        func: distribution to sample
    Returns:
        a list of exactly count things from the distribution
    """
    source = uniques(func(*vargs, **kwargs) for i in xrange(5 * count))
    return random.sample(list(source), k=count)


def generate_survey(choices, survey, comparisons, ratings, precision):
    """
    Generate all the Questions for a particular survey.

    Args:
        choices: things to rate, as list of HTML strings
        survey: survey URL
        comparisons: # of comparison questions
        ratings: # of rating questions
        precision: max # of stars, for rating questions
    Returns:
        list of Questions
    """

    return (
        repeat(comparisons, generate_comparison, choices, survey) +
        repeat(ratings, generate_rating, choices, survey, precision)
    )


def grouper(iterable, n, fillvalue=None):
    """
    Splits the iterable into groups of n elements.

    Args:
        iterable: things to process
        n: # of things per group
        fillvalue: what to do with the leftovers
    Returns:
        a list of groups
    """

    args = [iter(iterable)] * n
    return itertools.izip_longest(*args, fillvalue=fillvalue)


def uniques(values):
    """
    Remove any duplicates from the values.

    Args:
        values: list of things to search for
    Returns:
        a list of things, without duplicates
    """

    exist = {}
    for item in values:
        if item in exist:
            continue
        exist[item] = True
        yield item


def generate_surveys(choices, precisions=[2, 5, 100, "cmp"], count=100):
    """
    Generate a list of iterables for each survey.

    Args:
        choices: things to rate, as list of HTML strings
        precisions: list of rating precisions or "cmp" for only comparisons
        count: # of surveys
    Returns:
        a list of lists of Questions, one for each survey
    """

    surveys = []
    for prec in precisions:

        cmps, rates = int(count * 0.2), int(count * 0.8)
        if prec == "cmp":
            cmps, rates = count, 0

        for i in xrange(count):
            url = "experiment/{}/{}".format(i, prec)
            surveys.append(generate_survey(choices, url, cmps, rates, prec))

    return surveys


def main():
    meme_urls = open("memes.txt").readlines()
    markup = ['<!--{}--><img src="{}" style="width:90%"/>'.format(idx, url)
              for idx, url in enumerate(meme_urls)]

    shutil.rmtree("surveys", ignore_errors=True)
    os.mkdir("surveys")

    for i, surveys in enumerate(grouper(generate_surveys(markup), 30)):
        with open("surveys/chunk-{:02d}.json".format(i), "w") as fp:
            for survey in surveys:
                for question in survey or []:
                    fp.write(json.dumps(question._asdict()) + "\n")

if __name__ == "__main__":
    main()
