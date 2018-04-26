#!/usr/bin/env perl

use warnings;
use strict;

use constant COLUMNS => 2;
use constant COLUMN_WIDTH => 24;

my $list;

my $list_only;
if ($#ARGV >= 0 and $ARGV[0] =~ /^--?l(ist([-_]?only)?)?$/i) {
  shift;
  $list_only = 1;
}

while (<>) {
  if (/^\.(.*)/) {
    $list .= "$1\n";
  } else {
    if (length($list)) {
      $list = "" unless defined($list);
      if ($list_only) {
        print "$_\n" for sort split(/\s+/, $list);
      } else {
        my $columnated;
        my $column = 0;
        my $space = "";
        for my $id (sort split(/\s+/, $list)) {
          my $new_column = $column + 1 + int(length($id) / COLUMN_WIDTH);
          if ($new_column > COLUMNS) {
            $columnated .= "\n";
            $new_column -= $column;
            $column = 0;
            $space = "";
          }
          $columnated .= $space . $id;
          $space = " " x (($new_column-$column)*COLUMN_WIDTH - length($id));
          $space .= " " x ($new_column-$column-1);
          $column = $new_column;
        }
        $columnated =~ s/_/\\_/g;
        $columnated =~ s/\? /\?\\ /g;
        $columnated =~ s/((?:\\.|\S)+)/{\\cf $1}/g;
        print "$columnated\n";
        undef $list;
      }
    }
    print unless $list_only;
  }
}
