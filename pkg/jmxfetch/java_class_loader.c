#include <stdlib.h>
#include <stdio.h>
#include <stdbool.h>
#include "java_class_loader.h"

int run_class(jvm_loader* jloader, jclass java_class, char *method_name, char *method_signature)
{
	jmethodID method;

	method = (*jloader->env)->GetStaticMethodID(
			jloader->env,
			java_class,
			method_name,
			method_signature);
	if (method == NULL) {
		printf("Could not GetStaticMethodID '%s' with signature '%s'\n", method_name, method_signature);
		return false;
	}

	char *jmx_args[] = {
		"--ipc_host=localhost",
		"--ipc_port=5001",
		"--check_period=15000",
		"--thread_pool_size=3",
		"--collection_timeout=60",
		"--reconnection_timeout=10",
		"--reconnection_thread_pool_size=3",
		"--log_level=INFO",
		"--reporter=statsd:localhost:8125",
		"collect",
		NULL
	};

	// args
	jobjectArray applicationArgs = (*jloader->env)->NewObjectArray(jloader->env,
			sizeof(jmx_args)/sizeof(jmx_args[0]) -1,
			find_class(jloader, "java/lang/String"), NULL);

	for (int i = 0; jmx_args[i] != NULL; i++) {
		printf("jmx args: %d %s\n", i, jmx_args[i]);
		jstring jargs = (*jloader->env)->NewStringUTF(jloader->env, jmx_args[i]);
		(*jloader->env)->SetObjectArrayElement(jloader->env, applicationArgs, i, jargs);
	}

	(*jloader->env)->CallStaticVoidMethod(jloader->env, java_class, method, applicationArgs);
	return true;
}

jclass find_class(jvm_loader *jloader, char *class_name)
{
	jclass java_class = (*jloader->env)->FindClass(jloader->env, class_name);
	if (java_class == NULL) {
		printf("could not find class: %s\n", class_name); // TODO: raise error
		return NULL;
	}
	return java_class;
}

jvm_loader *create_vm(int options_len, char **options)
{
	jvm_loader *jloader = malloc(sizeof(jvm_loader));
	if (jloader == NULL) {
		printf("Could not allocate jvm_loader\n"); // TODO: raise error
		return NULL;
	}
	jloader->jvm = NULL;
	jloader->env = NULL;

	JavaVMOption jvm_options[options_len];

	for (int i = 0; i < options_len; i++) {
		printf("java option: %s\n", options[i]);
		jvm_options[0].optionString = options[i];
	}

	JavaVMInitArgs vm_args;
	vm_args.version = JNI_VERSION_1_6;
	vm_args.nOptions = options_len;
	vm_args.options = jvm_options;
	vm_args.ignoreUnrecognized = false;

	int jni_result = JNI_CreateJavaVM(&jloader->jvm, (void**)&jloader->env, &vm_args);
	if (jni_result != JNI_OK) {
		printf("Error creating a VM (error code %d)", jni_result); // TODO: raise error
		free(jloader);
		return NULL;
	}
	return jloader;
}

