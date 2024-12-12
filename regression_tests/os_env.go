package regressiontests

import "os"

func setEnvURLOrDefault(env string) string {
	v, ok := os.LookupEnv(env)
	if !ok {
		v = "http://localhost:3003"
	}
	return v
}

func setEnvXPrivOrDefault(env string) string {
	const adminXPriv = "xprv9s21ZrQH143K3CbJXirfrtpLvhT3Vgusdo8coBritQ3rcS7Jy7sxWhatuxG5h2y1Cqj8FKmPp69536gmjYRpfga2MJdsGyBsnB12E19CESK"

	v, ok := os.LookupEnv(env)
	if !ok {
		v = adminXPriv
	}
	return v
}
