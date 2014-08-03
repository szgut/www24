from django.db import models

class Score(models.Model):
    task = models.CharField(max_length=100)
    snapshot = models.IntegerField()
    team = models.CharField(max_length=100)
    value = models.FloatField(db_column='score')

    class Meta:
        db_table = 'score'
        unique_together = ('task', 'snapshot', 'team')
        index_together = [['task', 'snapshot', 'team'],]

    def __unicode__(self):
        return "%s/%s#%s" % (self.team, self.task, self.snapshot)
