from django.shortcuts import render
from score.models import Score
import collections

def big_score(small_score, avg):
    if avg == 0:
        return 0
    return 100 * small_score / avg
    

def home(request, template_name='scores.html'):
    scores = filter(lambda s: not s.team.startswith('.'), Score.objects.filter(snapshot=0))
    
    task_scores = collections.defaultdict(list)
    for score in scores:
        task_scores[score.task].append(score.value)
    task_avg = {task : sum(sorted(scores)[-3:])/3 for task, scores in task_scores.iteritems()}
    
    values = collections.defaultdict(lambda: collections.defaultdict(float))
    values_big = collections.defaultdict(lambda: collections.defaultdict(float))
    teams = collections.defaultdict(float)
    for score in scores:
        big = big_score(score.value, task_avg[score.task])
        values_big[score.team][score.task] = big
        values[score.team][score.task] = score.value
        teams[score.team] += big
    
    return render(request, template_name, {
        'tasks': sorted(set(map(lambda score: score.task, scores))),
        'teams': sorted(((score, team) for team, score in teams.iteritems()), reverse=True),
        'scores': values,
        'scores_big': values_big,
    })