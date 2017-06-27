
.PHONY: all clean cleaner checkindex

all: r7rs.pdf overview.pdf

clean:
	rm -f *~

cleaner: clean
	rm -f *.pdf *.aux *.log *.idx *.toc *.out

DIFF_SOURCES=basic.tex derive.tex example.tex expr.tex \
	lex.tex procs.tex prog.tex struct.tex

SOURCES=r7rs.tex $(DIFF_SOURCES) intro.tex \
	bib.tex syn.tex commands.tex first.tex notes.tex \
	repository.tex index.tex sem.tex stdmod-raw.tex \
	features.tex

# index.tex: r7rs.idx
# 	csi index.sch -e '(run)'

stdmod.tex: stdmod-raw.tex
	./genstdmod.pl < $< > $@

checkindex: stdmod-raw.tex index.tex
	./checkindex.sh

intro-ebook.tex: intro.tex
	sed 's/\\clearextrapart{\(.*\)}/\1/g' $< > $@

r7rs.pdf: $(SOURCES) stdmod.tex
	pdflatex $<
	pdflatex $<

overview.pdf: overview.tex overview-body.tex
	pdflatex $<
	pdflatex $<

r7rs-ebook.pdf: r7rs-ebook.tex intro-ebook.tex $(SOURCES) stdmod.tex
	pdflatex $<
	pdflatex $<

r5diff/%.tex: %.tex
	hg cat -r 1 $< > r5diff/old-$<
	latexdiff --type=UNDERLINE --subtype=SAFE r5diff/old-$< $< | \
	  perl -pe 's/\\section{\\DIF(add|del){([^{}]*)}}/\\section{\2}/; s/}\\ev  \\DIFadd{/\\ev/' | \
	  perl -pe 'BEGIN{undef $$/} s/\\ev([^%{]*)(%.*)?\n}/} \\ev \1\n/g' > $@
	rm -f r5diff/old-$<

r5diff/r7rs.pdf: $(DIFF_SOURCES:%=r5diff/%)
	cd r5diff && $(MAKE) r7rs.pdf
