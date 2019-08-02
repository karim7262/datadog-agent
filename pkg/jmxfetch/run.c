#include "java_class_loader.h"

int run_jmx(int ac, char **av)
{
	printf("creating VM (%d args)\n", ac-1);
	jvm_loader *jloader = create_vm(ac, av);
	if (jloader == NULL) {
		printf("C: creating JVM failed\n");
		return 1;
	}

	printf("getting  class\n");
	jclass reporter = find_class(jloader, "org/datadog/jmxfetch/reporter/StatsdReporter");
	if (reporter == NULL)
		return 1;

	printf("Injecting aggregator methods\n");
	if (!inject_aggregator_methods(jloader, reporter))
		return 1;

	printf("getting JMX class\n");
	jclass jmx = find_class(jloader, "org/datadog/jmxfetch/App");
	if (jmx == NULL)
		return 1;

	printf("running main method\n");
	if (!run_class(jloader, jmx, "main", "([Ljava/lang/String;)V"))
		return 1;

	printf("success !\n");
	return 0;
}
