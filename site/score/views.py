from django.shortcuts import render
from score.models import Score

def home(request, template_name='scores.html'):      
    return render(request, template_name, {
        'teams': Score.objects.order_by('-score'),
    })