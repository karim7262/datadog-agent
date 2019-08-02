#include "java_class_loader.h"

static void gauge(JNIEnv* env, jobject java_this, jstring metric_name, jint value)
{
	const char* name = (*env)->GetStringUTFChars(env, metric_name, NULL);
	printf("C: Gauge '%s': %ld\n", name, (long)value);
}

static void count(JNIEnv* env, jobject java_this, jstring metric_name, jint value)
{
	const char* name = (*env)->GetStringUTFChars(env, metric_name, NULL);
	printf("C: Count '%s': %ld\n", name, (long)value);
}

static void histogram(JNIEnv* env, jobject java_this, jstring metric_name, jint value)
{
	const char* name = (*env)->GetStringUTFChars(env, metric_name, NULL);
	printf("C: histogram '%s': %ld\n", name, (long)value);
}

static JNINativeMethod native_methods[] = {
	{ "gauge", "(Ljava/lang/String;I)V", (void*)gauge},
	{ "count", "(Ljava/lang/String;I)V", (void*)count},
	{ "histogram", "(Ljava/lang/String;I)V", (void*)histogram},
};

bool inject_aggregator_methods(jvm_loader *jloader, jclass java_class)
{
	int num_methods = sizeof(native_methods) / sizeof(native_methods[0]);

	printf("num_methods: %d\n", num_methods);
	if ((*jloader->env)->RegisterNatives(jloader->env, java_class, native_methods, num_methods) < 0) {
		printf("could not add methods\n");
		return false;
	}
	return true;
}
