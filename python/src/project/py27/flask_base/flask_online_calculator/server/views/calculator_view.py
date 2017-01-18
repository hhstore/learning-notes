
from flask import render_template, flash, redirect, session, url_for, request, g

from ..base import app
from ..forms.calculator_form import CalculatorForm


@app.route('/', methods=['GET', 'POST'])
def calculator():
    form = CalculatorForm()
    result = 0

    num1 = request.form.get("number1")
    num2 = request.form.get("number2")
    if num1 and num2:
        try:
            result = int(num1) + int(num2)
        except ValueError:
            flash("data type is invalid !")
            print "data type is invalid !"

    return render_template("calc.html", form=form, result=result)
