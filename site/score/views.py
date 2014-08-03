from django.shortcuts import render
from score.models import Score
import collections

def home(request, template_name='scores.html'):
    scores = Score.objects.filter(snapshot=0)
    values = collections.defaultdict(lambda: collections.defaultdict(float))
    teams = collections.defaultdict(float)
    for score in scores:
        values[score.team][score.task] = score.value
        teams[score.team] += score.value
    return render(request, template_name, {
        'tasks': set(map(lambda score: score.task, scores)),
        'teams': sorted(((score, team) for team, score in teams.iteritems()), reverse=True),
        'scores': values,
    })