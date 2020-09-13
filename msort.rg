(def (merge l1 l2)
  (if (empty? l1)
    l2
    (if (empty? l2)
      l1
      (if (< (first l1) (first l2))
        (cons (first l1) (merge (rest l1) l2))
        (cons (first l2) (merge l1 (rest l2)))))))

(def (split l)
  (let loop ((l l) (acc1 '()) (acc2 '()))
    (if (empty? l)
      (list acc1 acc2)
      (loop (rest l) acc2 (cons (first l) acc1)))))

(def (msort l)
  (if (or (empty? l) (empty? (rest l)))
    l
    (apply merge (map msort (split l)))))



;; With cond:

;; (def (merge l1 l2)
;;   (cond ((empty? l1) l2)
;;         ((empty? l2) l1)
;; 	((< (first l1) (first l2)) (cons (first l1) (merge (rest l1) l2)))
;; 	(cons (first l2) (merge l1 (rest l2)))))

;; (def (split l)
;;   (cond ((empty? l) '(() ()))
;;         ((empty? (rest l)) `((,(first l)) ()))
;; 	(let ((split-rest (split (rest (rest l)))))
;; 	  `(,(cons (first l) (first split-rest))
;; 	    ,(cons (second l) (second split-rest))))))

;; (def (msort l)
;;   (cond ((or (empty? l) (empty? (rest l))) l)
;;         (apply merge (map msort (split l)))))


;; With matching:

;; (def (merge l1 l2)
;;   (match (list l1 l2)
;;     ( (() _)                l2)
;;     ( (_ ())                l1)
;;     ( ((f1 . r1) (f2 . r2)) (if (< f1 g2)
;;                               (cons f1 (merge r1 l2))
;;                               (cons f2 (merge l1 r2))))))
			     
;; (def (split l)
;;   (match l
;;     ( ()          '(() ()))
;;     ( (f)         `((,f) ()))
;;     ( (f1 f2 . r) (match (split r)
;;                    ((r1 r2) (list (cons f1 r1) (cons f2 r2)))))))

;; (def (msort l)
;;   (match l
;;     ( ()  l)
;;     ( (_) l)
;;     ( _   (apply merge (map msort (split l))))))
    
