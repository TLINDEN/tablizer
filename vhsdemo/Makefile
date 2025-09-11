.PHONY: demo check clean-demo

VHS  = vhs

clean-demo:

%.gif: %.tape
	@echo "vhs $<"
	env PATH=..:$(PATH) vhs $<

check:
	ls -l ../tablizer

demo: check clean-demo demo.gif

