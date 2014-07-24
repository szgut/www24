from django.db import models

class Team(models.Model):
    name = models.CharField(primary_key=True, max_length=50)
    score = models.FloatField()
    
    def __unicode__(self):
        return self.name

class Snapshot(Team):
    round = models.IntegerField()