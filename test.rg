(def (square a) (* a a))

(def (test b)
  (let ((s (square b)))
    (print "The square of" b "is" s)))
    
(print "hello")
