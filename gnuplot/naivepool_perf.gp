set terminal png

set output output_path 

set style fill solid border -1

set autoscale

set title "Execute 1000 Tasks"
set ylabel 'time(nsec)'

plot input_path using 2:xtic(1) with histograms title "original"
