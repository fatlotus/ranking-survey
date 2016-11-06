# Ranking Survey

A Google Forms-like surveying platform, designed to make it easy for software
to ask people questions. Doing this requires uploading the code to App Engine,
generating survey JSON files, and uploading them to the site.

## Tweaking the questions

To change what questions are asked, download this repository and edit the
`make_questions.py` file. Next, rerun the application to generate the survey
JSON files in `surveys/`:

	$ python make_questions.py

Once you've run `make_questions.py`, You can see the options that can be changed
by poking around the JSON files.

	$ head -n1 surveys/chunk-00.json | python -mjson.tool
	{
		"choices": [
			"<!--24--><img src=\"http://i.giphy.com/E600ju0QaVvtC.gif\n\" style=\"width:90%\"/>",
			"<!--81--><img src=\"http://i.giphy.com/YB7l8W2j0yL96.gif\n\" style=\"width:90%\"/>"
		],
		"exclusive": true,
		"precision": 2,
		"survey": "experiment/0/2"
	}

In this case, `choices` are the HTML options that are presented to the user.
The `exclusive` option determines whether the user can give two options the
same level of preference (whether two things can be ranked highest). The
`precision` field indicates how specific the user can be.

Finally, `survey` indicates which URL this question will appear under.
Questions are only shown once, so usually one survey is generated for each
person.

You can see the specification of these records in `types.go`.

## Code changes

To run a local test instance of the application, download and install
[Go](https://golang.org/dl), set your GOPATH, then run:

	$ go get github.com/fatlotus/rankingsurvey

to download the code in this repository. To run a local test instance of the
application, run

	$ go install github.com/fatlotus/rankingsurvey/...
	$ $GOPATH/bin/rankings

and visit `http://localhost:8080` in a web browser. To sign in as an
administrator in, visit `http://localhost:8080/admin`, select the
"administrator" box, then visit `http://localhost:8080`. You should then be
able to upload the surveys you've generated earlier.

## Uploading the code to App Engine

Since the Github repository is automatically set to perform a CircleCI
deployment to App Engine, you should not ever need to do this.

Moreover, you only need to re-upload the application if the Go code changes;
changing the questions just requires using the administration page.

However, should the need arise, you will need to [download and
configure the gcloud command](https://cloud.google.com/sdk/gcloud/),
then run:

    $ gcloud -q app deploy app.yaml --promote --version=pre-gauss

That will build and depoy the code to app engine. If prompted, follow the
instructions in the console.

## License

The code in this repository is covered under the MIT License:

> Copyright (c) 2015 Jeremy Archer
> 
> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
> 
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
> 
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
> BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
> ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
> CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
> SOFTWARE.
