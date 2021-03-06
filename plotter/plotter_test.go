package plotter

import (
	"math"
	"testing"

	"github.com/go-test/deep"
)

func pS(s string) *string {
	return &s
}

func pF(f float64) *float64 {
	return &f
}

func pI(i int) *int {
	return &i
}

func TestParser(t *testing.T) {
	input := "-240 - i, 110 + 10 * sin(i / 5), -270 + 10 * cos(i / 5) | i <- 1..100"
	actual := parse(input)
	expected := function{
		Expressions: &expressions{
			X: &expression{
				Left: &term{
					Left: &factor{
						Base: &value{
							Number: pF(-240),
						},
					},
				},
				Right: []*opTerm{
					{
						Operator: opSub,
						Term: &term{
							Left: &factor{
								Base: &value{
									Ident: pS("i"),
								},
							},
						},
					},
				},
			},
			Y: &expression{
				Left: &term{
					Left: &factor{
						Base: &value{
							Number: pF(110),
						},
					},
				},
				Right: []*opTerm{
					{
						Operator: opAdd,
						Term: &term{
							Left: &factor{
								Base: &value{
									Number: pF(10),
								},
							},
							Right: []*opFactor{
								{
									Operator: opMul,
									Factor: &factor{
										Base: &value{
											Call: &call{
												Name: pS("sin"),
												Args: []*expression{
													&expression{
														Left: &term{
															Left: &factor{
																Base: &value{
																	Ident: pS("i"),
																},
															},
															Right: []*opFactor{
																{
																	Operator: opDiv,
																	Factor: &factor{
																		Base: &value{
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
			Z: &expression{
				Left: &term{
					Left: &factor{
						Base: &value{
							Number: pF(-270),
						},
					},
				},
				Right: []*opTerm{
					{
						Operator: opAdd,
						Term: &term{
							Left: &factor{
								Base: &value{
									Number: pF(10),
								},
							},
							Right: []*opFactor{
								{
									Operator: opMul,
									Factor: &factor{
										Base: &value{
											Call: &call{
												Name: pS("cos"),
												Args: []*expression{
													&expression{
														Left: &term{
															Left: &factor{
																Base: &value{
																	Ident: pS("i"),
																},
															},
															Right: []*opFactor{
																{
																	Operator: opDiv,
																	Factor: &factor{
																		Base: &value{
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
		Ranges: []*identRange{
			&identRange{Ident: pS("i"), Start: pI(1), End: pI(100)},
		},
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}
}

func TestEvalExpression(t *testing.T) {
	exp := expression{
		Left: &term{
			Left: &factor{
				Base: &value{
					Number: pF(110),
				},
			},
		},
		Right: []*opTerm{
			{
				Operator: opAdd,
				Term: &term{
					Left: &factor{
						Base: &value{
							Number: pF(10),
						},
					},
					Right: []*opFactor{
						{
							Operator: opMul,
							Factor: &factor{
								Base: &value{
									Call: &call{
										Name: pS("sin"),
										Args: []*expression{
											&expression{
												Left: &term{
													Left: &factor{
														Base: &value{
															Ident: pS("i"),
														},
													},
													Right: []*opFactor{
														{
															Operator: opDiv,
															Factor: &factor{
																Base: &value{
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
	}
	e := env{
		"i": 20,
		"sin": func(fs ...float64) float64 {
			return math.Sin(fs[0])
		},
	}

	actual := exp.eval(e)
	expected := 102.431975047
	if math.Abs(actual-expected) >= 0.000001 {
		t.Errorf("%f != %f", actual, expected)
	}
}

func TestMergeEnvs(t *testing.T) {
	actual := mergeEnvs(env{"i": 1}, env{"j": 2})
	expected := env{"i": 1, "j": 2}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}
}

func TestMergeListsOfEnvs(t *testing.T) {
	actual := mergeListsOfEnvs(
		[]env{
			{"i": 1},
			{"i": 2},
		},
		[]env{
			{"j": 1},
			{"j": 2},
		},
	)
	expected := []env{
		{"i": 1, "j": 1},
		{"i": 1, "j": 2},
		{"i": 2, "j": 1},
		{"i": 2, "j": 2},
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}
}

func TestEvalFunction(t *testing.T) {
	ast := function{
		Expressions: &expressions{
			X: &expression{
				Left: &term{
					Left: &factor{
						Base: &value{
							Ident: pS("i"),
						},
					},
				},
			},
			Y: &expression{
				Left: &term{
					Left: &factor{
						Base: &value{
							Ident: pS("j"),
						},
					},
				},
			},
			Z: &expression{
				Left: &term{
					Left: &factor{
						Base: &value{
							Ident: pS("k"),
						},
					},
				},
			},
		},
		Ranges: []*identRange{
			&identRange{Ident: pS("i"), Start: pI(1), End: pI(2)},
			&identRange{Ident: pS("j"), Start: pI(1), End: pI(2)},
			&identRange{Ident: pS("k"), Start: pI(1), End: pI(2)},
		},
	}
	actual := ast.eval(nil)
	expected := []coord{
		{1, 1, 1},
		{1, 1, 2},
		{1, 2, 1},
		{1, 2, 2},
		{2, 1, 1},
		{2, 1, 2},
		{2, 2, 1},
		{2, 2, 2},
	}
	if diff := deep.Equal(actual, expected); diff != nil {
		t.Error(diff)
	}
}

func TestOneArgPreludeFuncs(t *testing.T) {
	tests := []struct {
		fnName string
		result float64
	}{
		{"sin", -0.95892427466},
		{"cos", 0.28366218546},
	}
	for _, test := range tests {
		t.Run(test.fnName, func(t *testing.T) {
			exp := expression{
				Left: &term{
					Left: &factor{
						Base: &value{
							Call: &call{
								Name: pS(test.fnName),
								Args: []*expression{
									&expression{
										Left: &term{
											Left: &factor{
												Base: &value{
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
			}

			actual := exp.eval(prelude)
			expected := test.result
			if math.Abs(actual-expected) >= 0.000001 {
				t.Errorf("%f != %f", actual, expected)
			}
		})

	}
}
