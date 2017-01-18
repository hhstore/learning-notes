# -*- coding:utf-8  -*-

from flask_wtf import Form
from wtforms import TextField, BooleanField, PasswordField
from wtforms.validators import Required


class CalculatorForm(Form):
    number1 = TextField('Name', validators=[Required()])
    number2 = TextField('Name', validators=[Required()])
    choice = TextField('Name', validators=[Required()])
    result = TextField('Name', validators=[Required()])


