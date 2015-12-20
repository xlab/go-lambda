#ifdef _POSIX_C_SOURCE
#undef _POSIX_C_SOURCE
#endif

#include "Python.h"
#include "structmember.h"
#include "memoryobject.h"
#include "bufferobject.h"

#include "module.h"

#if PY_VERSION_HEX > 0x03000000
#error "python 3 is not yet supported"
#endif

typedef void* lambda_handler_func;

typedef struct {
	PyObject_HEAD
	lambda_handler_func ptr;
} lambda_handler;

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

static PyObject *
lambda_handler_tp_call(lambda_handler *self, PyObject *args, PyObject *other) {
	const char *event;
	const char *context;
	if (!PyArg_ParseTuple(args, "ss", &event, &context)) {
		return NULL;
	}
	lambda_handler_result r = lambda_handler_call(self->ptr, event, context);
	PyObject *result = PyString_FromStringAndSize(r.data, r.size);
	return result;
}

static PyTypeObject lambda_handler_type = {
	PyObject_HEAD_INIT(NULL)
	0,	/*ob_size*/
	"lambda_handler",		/*tp_name*/
	sizeof(lambda_handler),	/*tp_basicsize*/
	0,	/*tp_itemsize*/
	(destructor)lambda_handler_dealloc,	/*tp_dealloc*/
	0,	/*tp_print*/
	0,	/*tp_getattr*/
	0,	/*tp_setattr*/
	0,	/*tp_compare*/
	0,	/*tp_repr*/
	0,	/*tp_as_number*/
	0,	/*tp_as_sequence*/
	0,	/*tp_as_mapping*/
	0,	/*tp_hash */
	(ternaryfunc)lambda_handler_tp_call,	/*tp_call*/
	0,	/*tp_str*/
	0,	/*tp_getattro*/
	0,	/*tp_setattro*/
	0,	/*tp_as_buffer*/
	Py_TPFLAGS_DEFAULT,	/*tp_flags*/
	"",	/* tp_doc */
	0,	/* tp_traverse */
	0,	/* tp_clear */
	0,	/* tp_richcompare */
	0,	/* tp_weaklistoffset */
	0,	/* tp_iter */
	0,	/* tp_iternext */
	0,  /* tp_methods */
	0,	/* tp_members */
	0,	/* tp_getset */
	0,	/* tp_base */
	0,	/* tp_dict */
	0,	/* tp_descr_get */
	0,	/* tp_descr_set */
	0,	/* tp_dictoffset */
	0,  /* tp_init */
	0, 	/* tp_alloc */
	lambda_handler_new,	/* tp_new */
};

PyMODINIT_FUNC
init_module(void)
{
	if (PyType_Ready(&lambda_handler_type) < 0) { return; }
	PyObject *module = Py_InitModule3("go_lambda_handler", 0, "");

	Py_INCREF(&lambda_handler_type);
	PyModule_AddObject(module, "lambda_handler", (PyObject*)&lambda_handler_type);
}

