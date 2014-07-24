from django.db import models

class Score(models.Model):
    team = models.CharField(max_length=100)
    task = models.CharField(max_length=100)
    snapshot = models.IntegerField()
    score = models.FloatField()

    class Meta:
        db_table = 'score'
        unique_together = ('team', 'task', 'snapshot')

    def __unicode__(self):
        return "%s/%s#%s" % (self.team, self.task, self.snapshot)
