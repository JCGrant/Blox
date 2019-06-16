package plotter

import (
	"testing"

	"github.com/go-test/deep"
)

func pS(s string) *string {
	return &s
}

func pF(i float64) *float64 {
	return &i
}

func TestParser(t *testing.T) {
	input := "-240 - i, 110 + 10 * sin(i / 5), -270 + 10 * cos(i / 5) | i <- 1..100"
	actual := parse(input)
	expected := function{
		Expressions: expressions{
			X: expression{
				Left: term{
					Left: factor{
						Base: value{
							Number: pF(-240),
						},
					},
				},
				Right: []opTerm{
					{
						Operator: opSub,
						Term: term{
							Left: factor{
								Base: value{
									Ident: pS("i"),
								},
							},
						},
					},
				},
			},
			Y: expression{
				Left: term{
					Left: factor{
						Base: value{
							Number: pF(110),
						},
					},
				},
				Right: []opTerm{
					{
						Operator: opAdd,
						Term: term{
							Left: factor{
								Base: value{
									Number: pF(10),
								},
							},
							Right: []opFactor{
								{
									Operator: opMul,
									Factor: factor{
										Base: value{
											Call: &call{
												Name: "sin",
												Args: []expression{
													expression{
														Left: term{
															Left: factor{
																Base: value{
																	Ident: pS("i"),
																},
															},
															Right: []opFactor{
																{
																	Operator: opDiv,
																	Factor: factor{
																		Base: value{
																			Number: pF(5),
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Z: expression{
				Left: term{
					Left: factor{
						Base: value{
							Number: pF(-270),
						},
					},
				},
				Right: []opTerm{
					{
						Operator: opAdd,
						Term: term{
							Left: factor{
								Base: value{
									Number: pF(10),
								},
							},
							Right: []opFactor{
								{
									Operator: opMul,
									Factor: factor{
										Base: value{
											Call: &call{
												Name: "cos",
												Args: []expression{
													expression{
														Left: term{
															Left: factor{
																Base: value{
																	Ident: pS("i"),
																},
															},
															Right: []opFactor{
																{
																	Operator: opDiv,
																	Factor: factor{
																		Base: value{
																			Number: pF(5),
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		Ranges: []identRange{
			identRange{Ident: "i", Start: 1, End: 100},
		},
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}
}

func TestEvalExpression(t *testing.T) {
	exp := expression{
		Left: term{
			Left: factor{
				Base: value{
					Number: pF(-240),
				},
			},
		},
		Right: []opTerm{
			{
				Operator: opSub,
				Term: term{
					Left: factor{
						Base: value{
							Ident: pS("i"),
						},
					},
				},
			},
		},
	}
	e := env{
		"i": 20,
	}

	actual := exp.eval(e)
	expected := -260.0
	if actual != expected {
		t.Errorf("%f != %f", actual, expected)
	}
}
