from django.shortcuts import render
from score.models import Team

def home(request, template_name='scores.html'):      
    return render(request, template_name, {
        'teams': Team.objects.order_by('-score'),
    })