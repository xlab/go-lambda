#ifdef _POSIX_C_SOURCE
#undef _POSIX_C_SOURCE
#endif

#include "Python.h"
#include "structmember.h"
#include "memoryobject.h"
#include "bufferobject.h"

#if PY_VERSION_HEX > 0x03000000
#error "python 3 is not yet supported"
#endif

typedef void* lambda_handler_func;

typedef struct {
	PyObject_HEAD
	lambda_handler_func ptr;
} lambda_handler;

extern void* get_lambda_handler();

static PyObject*
lambda_handler_new(PyTypeObject *type, PyObject *args, PyObject *kwds) {
	lambda_handler *self;
	self = (lambda_handler *)type->tp_alloc(type, 0);
	self->ptr = get_lambda_handler();
	return (PyObject*)self;
}

static void
lambda_handler_dealloc(lambda_handler *self) {
	self->ob_type->tp_free((PyObject*)self);
}

struct lambda_handler_call_return {
	char* data;
	size_t size;
};

extern struct lambda_handler_call_return lambda_handler_call(void* p0, char* event, char* context);

static PyObject *
lambda_handler_tp_call(lambda_handler *self, PyObject *args, PyObject *other) {
	char *event;
	char *context;
	if (!PyArg_ParseTuple(args, "ss", &event, &context)) {
		return NULL;
	}
	struct lambda_handler_call_return ret;
	ret = lambda_handler_call(self->ptr, event, context);
	PyObject *result = PyString_FromStringAndSize(ret.data, ret.size);
	return result;
}

static PyTypeObject lambda_handler_type = {
	PyObject_HEAD_INIT(NULL)
	0,                         /*ob_size*/
	"lambda_handler",          /*tp_name*/
	sizeof(lambda_handler),    /*tp_basicsize*/
	0,                         /*tp_itemsize*/
	(destructor)lambda_handler_dealloc,	/*tp_dealloc*/
	0,                         /*tp_print*/
	0,                         /*tp_getattr*/
	0,                         /*tp_setattr*/
	0,                         /*tp_compare*/
	0,                         /*tp_repr*/
	0,                         /*tp_as_number*/
	0,                         /*tp_as_sequence*/
	0,                         /*tp_as_mapping*/
	0,                         /*tp_hash */
	(ternaryfunc)lambda_handler_tp_call,	/*tp_call*/
	0,                         /*tp_str*/
	0,                         /*tp_getattro*/
	0,                         /*tp_setattro*/
	0,                         /*tp_as_buffer*/
	Py_TPFLAGS_DEFAULT,        /*tp_flags*/
	"",           			   /* tp_doc */
};

PyMODINIT_FUNC
initmodule(void)
{
	PyObject *module;

	lambda_handler_type.tp_new = lambda_handler_new;
	if (PyType_Ready(&lambda_handler_type) < 0) {
		return;
	} else {
		Py_INCREF(&lambda_handler_type);
	}

	module = Py_InitModule("module", NULL);
	PyModule_AddObject(module, "lambda_handler", (PyObject*)&lambda_handler_type);
}

