# Requiremets for the system
Installed:
- docker (with ability to run from the user)
- make

# Build
```sh
$ make build
```
<details><summary>Example</summary>

```sh
$ make build
docker build \
        --build-arg GO_VER=1.18.3 \
        --build-arg ALPINE_VER=3.15 \
        --build-arg WORKDIR=/go/src/github.com/hellofreshdevtests/recipe-count-task \
        -t vasily.chertkov/recipe-count-task:1.0.0 -f /home/vchertkov/dev/hellofresh/vasily-chertkov-recipe-count-test-2020/docker/Dockerfile .
Sending build context to Docker daemon    489kB
Step 1/12 : ARG GO_VER
Step 2/12 : ARG ALPINE_VER
Step 3/12 : FROM golang:${GO_VER}-alpine${ALPINE_VER} as builder
 ---> 391370347e38
Step 4/12 : LABEL stage=server-intermediate
 ---> Using cache
 ---> 5343332a5b2b
Step 5/12 : ARG WORKDIR
 ---> Using cache
 ---> ae798080bcae
Step 6/12 : WORKDIR ${WORKDIR}
 ---> Using cache
 ---> becefe459f5b
Step 7/12 : RUN apk --no-cache --update add     git
 ---> Using cache
 ---> f932cce0a026
Step 8/12 : COPY ./ ./
 ---> 7577f95770dc
Step 9/12 : RUN CGO_ENABLED=0 GOOS=linux go build -v -mod=vendor -o /tmp/recipe-count-task ./cmd/cli/*.go
 ---> Running in 3f92601bd059
github.com/modern-go/reflect2
github.com/modern-go/concurrent
github.com/json-iterator/go
command-line-arguments
Removing intermediate container 3f92601bd059
 ---> 83a170d0b636
Step 10/12 : FROM alpine:${ALPINE_VER} as base
 ---> c059bfaa849c
Step 11/12 : COPY --from=builder /tmp/recipe-count-task /usr/local/bin/recipe-count-task
 ---> 5a58953be0dc
Step 12/12 : ENTRYPOINT ["/usr/local/bin/recipe-count-task"]
 ---> Running in 3026ac3f573f
Removing intermediate container 3026ac3f573f
 ---> e603069252c3
Successfully built e603069252c3
Successfully tagged vasily.chertkov/recipe-count-task:1.0.0
```
</details>

# Run
The service requires the json file with the data to be passed.
It can be done by passsing the path to the built container.
The input file has to be mounted to the docker container as `/input.json` (See example below).
The filter flags:
- `from` - starting time of the delivery
- `to` - end time of the delivery
- `postcode` - postcode for the filtering deliveries
- `f` - filter for the recipes match, allowed multiple values 

```sh
$ docker run --rm -v /tmp/hf_test_calculation_fixtures.json:/input.json vasily.chertkov/recipe-count-task:1.0.0 -from 11AM -to 4PM -postcode 10120 -f Mex -f Yellow -f Mac

```
<details><summary>Example</summary>

```sh
{
    "unique_recipe_count": 29,
    "count_per_recipe": [
    {
        "recipe": "Cajun-Spiced Pulled Pork",
        "count": 667365
    },
    {
        "recipe": "Cheesy Chicken Enchilada Bake",
        "count": 333012
    },
    {
        "recipe": "Cherry Balsamic Pork Chops",
        "count": 333889
    },
    {
        "recipe": "Chicken Pineapple Quesadillas",
        "count": 331514
    },
    {
        "recipe": "Chicken Sausage Pizzas",
        "count": 333306
    },
    {
        "recipe": "Creamy Dill Chicken",
        "count": 333103
    },
    {
        "recipe": "Creamy Shrimp Tagliatelle",
        "count": 333395
    },
    {
        "recipe": "Crispy Cheddar Frico Cheeseburgers",
        "count": 333251
    },
    {
        "recipe": "Garden Quesadillas",
        "count": 333284
    },
    {
        "recipe": "Garlic Herb Butter Steak",
        "count": 333649
    },
    {
        "recipe": "Grilled Cheese and Veggie Jumble",
        "count": 333742
    },
    {
        "recipe": "Hearty Pork Chili",
        "count": 333355
    },
    {
        "recipe": "Honey Sesame Chicken",
        "count": 333748
    },
    {
        "recipe": "Hot Honey Barbecue Chicken Legs",
        "count": 334409
    },
    {
        "recipe": "Korean-Style Chicken Thighs",
        "count": 333069
    },
    {
        "recipe": "Meatloaf à La Mom",
        "count": 333570
    },
    {
        "recipe": "Mediterranean Baked Veggies",
        "count": 332939
    },
    {
        "recipe": "Melty Monterey Jack Burgers",
        "count": 333264
    },
    {
        "recipe": "Mole-Spiced Beef Tacos",
        "count": 332993
    },
    {
        "recipe": "One-Pan Orzo Italiano",
        "count": 333109
    },
    {
        "recipe": "Parmesan-Crusted Pork Tenderloin",
        "count": 333311
    },
    {
        "recipe": "Spanish One-Pan Chicken",
        "count": 333291
    },
    {
        "recipe": "Speedy Steak Fajitas",
        "count": 333578
    },
    {
        "recipe": "Spinach Artichoke Pasta Bake",
        "count": 333545
    },
    {
        "recipe": "Steakhouse-Style New York Strip",
        "count": 333473
    },
    {
        "recipe": "Stovetop Mac 'N' Cheese",
        "count": 333098
    },
    {
        "recipe": "Sweet Apple Pork Tenderloin",
        "count": 332595
    },
    {
        "recipe": "Tex-Mex Tilapia",
        "count": 333749
    },
    {
        "recipe": "Yellow Squash Flatbreads",
        "count": 333394
    }
],
    "busiest_postcode": {
    "postcode": "10176",
    "delivery_count": 91785
},
    "count_per_postcode_and_time": {
    "postcode": "10120",
    "delivery_count": 2939,
    "from": "11AM",
    "to": "4PM"
},
    "match_by_name": [
    "Stovetop Mac 'N' Cheese",
    "Tex-Mex Tilapia",
    "Yellow Squash Flatbreads"
]
}
Processing time 11.075531875s
```
</details>

Recipe Stats Calculator
====

In the given assignment we suggest you to process an automatically generated JSON file with recipe data and calculated some stats.

Instructions
-----

1. Clone this repository.
2. Create a new branch called `dev`.
3. Create a pull request from your `dev` branch to the master branch.
4. Reply to the thread you're having with your recruiter telling them we can start reviewing your code

Given
-----

Json fixtures file with recipe data. Download [Link](https://test-golang-recipes.s3-eu-west-1.amazonaws.com/recipe-calculation-test-fixtures/hf_test_calculation_fixtures.tar.gz)

_Important notes_

1. Property value `"delivery"` always has the following format: "{weekday} {h}AM - {h}PM", i.e. "Monday 9AM - 5PM"
2. The number of distinct postcodes is lower than `1M`, one postcode is not longer than `10` chars.
3. The number of distinct recipe names is lower than `2K`, one recipe name is not longer than `100` chars.

Functional Requirements
------

1. Count the number of unique recipe names.
2. Count the number of occurences for each unique recipe name (alphabetically ordered by recipe name).
3. Find the postcode with most delivered recipes.
4. Count the number of deliveries to postcode `10120` that lie within the delivery time between `10AM` and `3PM`, examples _(`12AM` denotes midnight)_:
    - `NO` - `9AM - 2PM`
    - `YES` - `10AM - 2PM`
5. List the recipe names (alphabetically ordered) that contain in their name one of the following words:
    - Potato
    - Veggie
    - Mushroom

Non-functional Requirements
--------

1. The application is packaged with [Docker](https://www.docker.com/).
2. Setup scripts are provided.
3. The submission is provided as a `CLI` application.
4. The expected output is rendered to `stdout`. Make sure to render only the final `json`. If you need to print additional info or debug, pipe it to `stderr`.
5. It should be possible to (implementation is up to you):  
    a. provide a custom fixtures file as input  
    b. provide custom recipe names to search by (functional reqs. 5)  
    c. provide custom postcode and time window for search (functional reqs. 4)  

Expected output
---------------

Generate a JSON file of the following format:

```json5
{
    "unique_recipe_count": 15,
    "count_per_recipe": [
        {
            "recipe": "Mediterranean Baked Veggies",
            "count": 1
        },
        {
            "recipe": "Speedy Steak Fajitas",
            "count": 1
        },
        {
            "recipe": "Tex-Mex Tilapia",
            "count": 3
        }
    ],
    "busiest_postcode": {
        "postcode": "10120",
        "delivery_count": 1000
    },
    "count_per_postcode_and_time": {
        "postcode": "10120",
        "from": "11AM",
        "to": "3PM",
        "delivery_count": 500
    },
    "match_by_name": [
        "Mediterranean Baked Veggies", "Speedy Steak Fajitas", "Tex-Mex Tilapia"
    ]
}
```

Review Criteria
---

We expect that the assignment will not take more than 3 - 4 hours of work. In our judgement we rely on common sense
and do not expect production ready code. We are rather instrested in your problem solving skills and command of the programming language that you chose.

It worth mentioning that we will be testing your submission against different input data sets.

__General criteria from most important to less important__:

1. Functional and non-functional requirements are met.
2. Prefer application efficiency over code organisation complexity.
3. Code is readable and comprehensible. Setup instructions and run instructions are provided.
4. Tests are showcased (_no need to cover everything_).
5. Supporting notes on taken decisions and further clarifications are welcome.

