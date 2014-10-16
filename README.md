track
=====

| This is a work in progress.

Command line tool to track how much time you spend on different projects.

To do (in big lines)
--------------------

- Finish all the basic commands
- Polish the syntax for those commands and make sure they're consistent
- Write tests
- Create a Visualisation
- Better structure the code

The first basic usable tool will come in the next few weeks and will contain
all the basic features (add, rm, init, delete, status, overview).

Usage
-----

Initialise index.

    $ track init

Add project names.

    $ track add <project1> <project2> ... <projectN>

Start working on a project. Use start or work:

    $ track work projectName
    $ track start projectName

Stop working.

    $ track stop

And other commands like delete (delete index), rm (remove projects), status (see current status).

Visualisation
-------------

Very basic but informative command line visualisation is in plan but the main way to visualise this data will
be through a web browser. The plan is to make a temporary HTML generator or a statis visualisation HTML page.

Getting the data to JavaScript from this tool is straightforward as it uses a very simple JSON file.
