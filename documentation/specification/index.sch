;; Program to process r7rs.idx entries.

;; Script for running this in Chicken:
;; $ csi index.sch -e '(run)'

(define (run)
  (call-with-input-file "r7rs.idx" read-entries)
  (call-with-output-file "index.tex" create-index))

;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(define main 0)
(define aux 1)

(define (make-entry key font main/aux page)
  (list key font main/aux page))
(define (entry-key x) (car x))
(define (entry-font x) (cadr x))
(define (entry-main/aux x) (caddr x))
(define (entry-page x) (cadddr x))

(define *database* '())

;; The TeX index entries for these values confuse the reader so we
;; skip them.
(define *bad-entries* '("" "##"))
(define *rename-entries*
  '(("fold-case@\\texttt  {\\#!fold-case}" . "\\#!fold-case")
    ("no-fold-case@\\texttt  {\\#!no-fold-case}" . "\\#!no-fold-case")))

(define (read-entries in)
  (let ((next (read-char in)))
    (cond ((eof-object? next))
	  ((char=? next #\")
	   (let ((key (read-rest-of-string in)))
             (cond
              ((member key *bad-entries*)
               (read-rest-of-line in)
               (read-entries in))
              (else
               (let* ((key (cond ((assoc key *rename-entries*) => cdr)
                                 (else key)))
                      (font (read in))
                      (main/aux (if (eq? 'main (read in)) 0 1))
                      (page (read in)))
                 (index-entry key font main/aux page)
                 (read-entries in))))))
	  (else
	   (read-entries in)))))

; For reading in strings that contain \.

(define (read-to-delimiter in delimiter)
  (let loop ((chars '()))
    (let ((next (read-char in)))
      (if (char=? next delimiter)
	  (list->string (reverse chars))
	  (loop (cons next chars))))))

(define (read-rest-of-string in)
  (read-to-delimiter in #\"))

(define (read-rest-of-line in)
  (read-to-delimiter in #\newline))

(define (index-entry key font main/aux page)
  (set! *database*
        (cons (make-entry (string-downcase key)
                          font
                          main/aux
                          page)
              *database*))
  #t)

(define (string-downcase string)
    (list->string (map char-downcase (string->list string))))

(define (create-index p)
  (define (loop)
    (if (null? *database*)
        'done
        (begin (process-key (collect-entries) p)
               (loop))))
  (set! *database*
        (sort *database*
              (lambda (x y)
                (string<? (entry-key x)
                          (entry-key y)))))
  (loop))

(define (collect-entries)
  (define (loop key entries)
    (cond ((null? *database*) entries)
          ((string=? key (entry-key (car *database*)))
           (let ((x (car *database*)))
             (set! *database* (cdr *database*))
             (loop key (cons x entries))))
          (else entries)))
  (loop (caar *database*) '()))

(define (process-key entries p)
  (let ((entries (sort entries entry<?)))
    (if (not (consistent? entries))
        (begin (display "Inconsistent entries:")
               (newline)
               (pretty-print entries)
               (newline)))
    (let ((key (entry-key (car entries)))
          (font (if (any? (lambda (e) (equal? "tt" (entry-font e)))
                         entries)
                    "tt"
                    (entry-font (car entries))))
          (main? (entry-main/aux (car entries)))
          (pages (remove-duplicates (map entry-page entries))))
      (if main?
          (write-entries key font (car pages) (cdr pages) p)
          (write-entries key font #f pages p)))))

(define (entry<? x y)
  (let ((x1 (entry-main/aux x))
        (y1 (entry-main/aux y)))
    (or (< x1 y1)
        (and (eq? x1 y1)
             (< (entry-page x) (entry-page y))))))

(define (consistent? entries)
  (let ((x (car entries)))
    (let ((key (entry-key x))
          (font (entry-font x)))
      (every? (lambda (x)
                (and (string=? key (entry-key x))
                     (or (string=? font (entry-font x))
			 ;; different entries for t and #t aren't inconsistent
			 (string=? font "sharpfoo")
			 (string=? (entry-font x) "sharpfoo"))
                     ;(eq? aux (entry-main/aux x))
                     ))
              (cdr entries)))))

(define (every? pred list)
    (let loop ((list list))
      (cond ((null? list)
             #t)
	    ((pred (car list))
	     (loop (cdr list)))
	    (else
	     #f))))

(define (any? pred list)
  (let loop ((list list))
    (cond ((null? list) #f)
          ((pred (car list)) #t)
          (else (loop (cdr list))))))

(define (remove-duplicates x)
  (define (loop x y)
    (cond ((null? x) (reverse y))
          ((memq (car x) y) (loop (cdr x) y))
          (else (loop (cdr x) (cons (car x) y)))))
  (loop (cdr x) (list (car x))))

(define *last-key* "%")
(define *s1* (string-append "\\item{" (list->string '(#\\))))
(define *s2* "{")
(define *s3* "}}{\\hskip .75em}")
(define *semi* "; ")
(define *comma* ", ")

(define (escape str)
  (if (equal? str "_")
      "\\_"
      str))

(define (write-page-ref page-n out)
  (display "\\hyperlink{page." out)
  (write page-n out)
  (display "}{" out)
  (write page-n out)
  (display "}" out))

(define (write-entries key font main pages p)
  (if (and (char-alphabetic? (string-ref key 0))
           (not (char=? (string-ref *last-key* 0)
                        (string-ref key 0))))
      (begin (display "\\indexspace" p)
             (newline p)))
  (set! *last-key* key)
  (display (string-append *s1* font *s2* (escape key) *s3*) p)
  (if main
      (begin (write-page-ref main p)
             (if (not (null? pages))
                 (display *semi* p))))
  (if (not (null? pages))
      (begin (write-page-ref (car pages) p)
             (for-each (lambda (page)
                         (display *comma* p)
                         (write-page-ref page p))
                       (cdr pages))))
  (newline p))
