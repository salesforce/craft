# Running mdBook
Docs are served using mdBook. If you want to test changes to the docs locally, follow these directions:

* Follow the instructions at https://github.com/rust-lang-nursery/mdBook#installation to install mdBook.
* Run mdbook serve
* Visit http://localhost:3000

# Steps to Deploy
Currently we are using Github Pages to deploy our docs. To send changes to the docs
* Run `mdbook build -d /tmp/craft-doc`
* `mv /tmp/craft-doc .`
* Send the Pull Request to `gh-pages`
