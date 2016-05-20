#!/usr/bin/env dtrace -qs

dtrace:::BEGIN
{
  printf("%-5s %-12s %-5s %-12s %6s\n","FROM", "", "TO", "", "SIGNAL");
}

proc:::signal-send
{
  printf("%5d %-12s %5d %-12s %6d\n",pid,execname,args[1]->pr_pid,args[1]->pr_fname,args[2]);
}
