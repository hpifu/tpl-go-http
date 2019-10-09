#!/usr/bin/env python3

from behave import *
from hamcrest import *
import json


@given('redis set object "{key:str}"')
def step_impl(context, key):
    obj = json.loads(context.text)
    context.redis_client.set(key, json.dumps(obj))


@given('redis set string "{key:str}"')
def step_impl(context, key):
    context.redis_client.set(key, context.text.strip())


@given('redis del "{key:str}"')
def step_impl(context, key):
    context.redis_client.delete(key)


@then('redis get object "{key:str}"')
def step_impl(context, key):
    obj = json.loads(context.text)
    val = context.redis_client.get(key)
    print(val)
    result = json.loads(val)
    for key in obj:
        assert_that(result[key], equal_to(obj[key]))


@then('redis not exist "{key:str}"')
def step_impl(context, key):
    res = context.redis_client.get(key)
    assert_that(res, equal_to(None))
