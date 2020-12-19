# ShakeSearch

Welcome to the Pulley Shakesearch Take-home Challenge! In this repository,
you'll find a simple web app that allows a user to search for a text string in
the complete works of Shakespeare.

You can see a live version of the app at
https://pulley-shakesearch.herokuapp.com/. Try searching for "Hamlet" to display
a set of results.

In it's current state, however, the app is just a rough prototype. The search is
case sensitive, the results are difficult to read, and the search is limited to
exact matches.

## Your Mission

Improve the search backend. Think about the problem from the user's perspective
and prioritize your changes according to what you think is most useful.

## Submission

1. Fork this repository and send us a link to your fork after pushing your changes. 
2. Heroku hosting - The project includes a Heroku Procfile and, in its
current state, can be deployed easily on Heroku's free tier.
3. In your submission, share with us what changes you made and how you would prioritize changes if you had more time.

## Tasks list - change log

- [ ] Highlight query text in result: like this is **query** string.
- [ ] Keeps formating (newlines) as in original text
- [ ] Provide phrases instead of [-250+i, i+250] text blocks
- [ ] Implement case insensitive search
- [ ] Avoid phrase repetition
- [ ] Implement phrase search: like `"love to hate you"`
- [ ] Allow several words in query: like `love hate`
- [ ] Return book title near the result phrase
- [ ] Prioritize search results based on statistics
