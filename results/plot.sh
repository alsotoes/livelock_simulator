echo "
  set datafile separator ','
  set ylabel 'Output packet rate'
  set xlabel 'Input packet rate'
  set title '$1'
  set grid
  set style circle radius screen 0.0001
  plot '$1' using 1:2 with circles fillstyle solid noborder lc rgb 'forest-green'
" | gnuplot --persist
