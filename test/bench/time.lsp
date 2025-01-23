; https://github.com/hymkor/lispect

(defglobal pid (spawn "go" "run" "slow.go"))
(expect ">")
(defglobal start (get-internal-real-time))
(sendln "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890")
(wait pid)
(defglobal end (get-internal-real-time))
(format (standard-output) "TIME: ~a~%" (div (convert (- end start) <float>) (internal-time-units-per-second)))
