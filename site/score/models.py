from django.db import models

class Team(models.Model):
    name = models.CharField(max_length=50)
    score = models.FloatField()
    
    def __unicode__(self):
        return self.name
