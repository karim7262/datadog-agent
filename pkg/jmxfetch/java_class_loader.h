#ifndef JAVA_CLASS_LOADER_H
#define JAVA_CLASS_LOADER_H

#include <stdbool.h>
#include <jni.h>

typedef struct jvm_loader_s {
	JNIEnv *env;
	JavaVM *jvm;
} jvm_loader;

// Load a class base on it's name
jclass find_class(jvm_loader *jloader, char *class_name);
// Create a new VM
jvm_loader *create_vm(int options_len, char **options);
// Inject aggregator methods in a Java class
bool inject_aggregator_methods(jvm_loader *jloader, jclass java_class);
// Call a method from a Java class
int run_class(jvm_loader* jloader, jclass java_class, char *method_name, char *method_signature);
// run_jmx
int run_jmx(int ac, char **av);

#endif // JAVA_CLASS_LOADER_H
