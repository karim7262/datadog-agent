// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.
#include "kubeutil.h"

#include "cgo_free.h"
#include <sixstrings.h>

// these must be set by the Agent
static cb_get_connection_info_t cb_get_connection_info = NULL;
static cb_collect_events_t cb_collect_events = NULL;

// forward declarations
static PyObject *get_connection_info();
static PyObject *collect_events();

static PyMethodDef methods[] = {
    { "get_connection_info", (PyCFunction)get_connection_info, METH_NOARGS, "Get kubelet connection information." },
    { "collect_events", (PyCFunction)collect_events, METH_NOARGS, "runs the event collection"},
    { NULL, NULL } // guards
};

#ifdef DATADOG_AGENT_THREE
static struct PyModuleDef module_def = { PyModuleDef_HEAD_INIT, KUBEUTIL_MODULE_NAME, NULL, -1, methods };

PyMODINIT_FUNC PyInit_kubeutil(void)
{
    return PyModule_Create(&module_def);
}
#endif

#ifdef DATADOG_AGENT_TWO
// in Python2 keep the object alive for the program lifetime
static PyObject *module;

void Py2_init_kubeutil()
{
    module = Py_InitModule(KUBEUTIL_MODULE_NAME, methods);
}
#endif

void _set_get_connection_info_cb(cb_get_connection_info_t cb)
{
    cb_get_connection_info = cb;
}

void _set_collect_events_cb(cb_collect_events_t cb)
{
    cb_collect_events = cb;
}

PyObject *collect_events(PyObject *self, PyObject *args)
{
    // callback must be set
    if (cb_collect_events == NULL)
        Py_RETURN_NONE;

    char *data;
    cb_collect_events(&data);

    // create a new ref
    PyObject *res_ver = PyStringFromCString(data);

    // free the memory allocated by the Agent
    cgo_free(data);

    if (res_ver == NULL ) {
        return PyDict_New();
    }
    return res_ver;
}

PyObject *get_connection_info(PyObject *self, PyObject *args)
{
    // callback must be set
    if (cb_get_connection_info == NULL)
        Py_RETURN_NONE;

    char *data = NULL;
    cb_get_connection_info(&data);

    // create a new ref
    PyObject *conn_info_dict = from_json(data);

    // free the memory allocated by the Agent
    cgo_free(data);

    if (conn_info_dict == NULL || !PyDict_Check(conn_info_dict)) {
        return PyDict_New();
    }
    return conn_info_dict;
}
