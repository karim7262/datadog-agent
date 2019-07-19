// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2019 Datadog, Inc.
#include "rtloader_mem.h"

#include <execinfo.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

// default memory management functions
static rtloader_malloc_t rt_malloc = malloc;
static rtloader_free_t rt_free = free;

// these must be set by the Agent
static cb_memory_tracker_t cb_memory_tracker = NULL;

void _set_memory_tracker_cb(cb_memory_tracker_t cb) {
    cb_memory_tracker = cb;
    srand(time(0));
}

static int my_traceback(char **traceback) {
#define MEM_TRACEBACK_SZ 128
#define SAMPLE_RATION 10
    int frames = 0;
    void *callstack[MEM_TRACEBACK_SZ];

    frames = backtrace(callstack, MEM_TRACEBACK_SZ);
    traceback = backtrace_symbols(callstack, frames);

    return frames;
}
void *_malloc(size_t sz) {
    int frames = 0;
    char **traceback = NULL;
    void *ptr = NULL;

    ptr = rt_malloc(sz);
    if (ptr && cb_memory_tracker) {
        if (!(rand() % SAMPLE_RATIO)) {
            frames = my_traceback(traceback);
        }
        cb_memory_tracker(ptr, sz, DATADOG_AGENT_RTLOADER_ALLOCATION, traceback, frames);

        if (traceback != NULL) {
            rt_free(traceback);
        }
    }

    return ptr;
}

void _free(void *ptr) {
    int frames = 0;
    char **traceback = NULL;

    rt_free(ptr);
    if (ptr && cb_memory_tracker) {
        if (!(rand() % SAMPLE_RATIO)) {
            frames = my_traceback(traceback);
        }
        cb_memory_tracker(ptr, 0, DATADOG_AGENT_RTLOADER_FREE, traceback, frames);

        if (traceback != NULL) {
            rt_free(traceback);
        }
    }
}

char *strdupe(const char *s1) {
    char * s2 = NULL;

    if (!(s2 = (char *)_malloc(strlen(s1)+1))) {
        return NULL;
    }

    return strcpy(s2, s1);
}
