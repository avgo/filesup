package main

/*

#cgo LDFLAGS: -lexif

#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <libexif/exif-loader.h>

#define TAG_VALUE_BUF    1024

struct tagExifResult {
	char DateTime[TAG_VALUE_BUF];
	char DateTimeOriginal[TAG_VALUE_BUF];
	char DateTimeDigitized[TAG_VALUE_BUF];
	char Make[TAG_VALUE_BUF];
	char Model[TAG_VALUE_BUF];
};

typedef struct tagExifResult ExifResult;

struct tagExifResultFields {
	ExifTag tag;
	int offset;
};

typedef struct tagExifResultFields ExifResultFields;

ExifResultFields erfd[] = {
	{ EXIF_TAG_DATE_TIME,            offsetof(struct tagExifResult, DateTime)          },
	{ EXIF_TAG_DATE_TIME_ORIGINAL,   offsetof(struct tagExifResult, DateTimeOriginal)  },
	{ EXIF_TAG_DATE_TIME_DIGITIZED,  offsetof(struct tagExifResult, DateTimeDigitized) },
	{ EXIF_TAG_MAKE,                 offsetof(struct tagExifResult, Make)              },
	{ EXIF_TAG_MODEL,                offsetof(struct tagExifResult, Model)             },
};

#define erfd_count (sizeof(erfd)/sizeof(*erfd))

int print_exif(const char *filename, ExifResult *result);
static void show1(ExifContent *c, void *data);
static void show2(ExifEntry *e, void *data);

int print_exif(const char *filename, ExifResult *result)
{
	int ret = 1;

	ExifLoader *l = exif_loader_new();

	exif_loader_write_file(l, filename);

	ExifData *ed = exif_loader_get_data(l);

	exif_loader_unref(l);

	if (ed == NULL)
	{
		fprintf(stderr, "error: some error occurs\n");
		ret = -1; goto END;
	}

	exif_data_foreach_content(ed, show1, result);

END:	exif_data_unref(ed);

	return ret;
}

static void show1(ExifContent *c, void *data)
{
	exif_content_foreach_entry(c, show2, data);
}

#define MY_ER(e)    ((ExifResult*)(e))

static void show2(ExifEntry *e, void *data)
{
	for (int i = 0; i < erfd_count; ++i)
	{
		if (erfd[i].tag == e->tag)
		{
			char *dst_buf = data + erfd[i].offset;
			exif_entry_get_value(e, dst_buf, TAG_VALUE_BUF);
			dst_buf[TAG_VALUE_BUF-1] = '\0';
		}
	}
}

*/
import "C"

import (
	"unsafe"
)

type ExifResult struct {
	DateTime string
	DateTimeOriginal string
	DateTimeDigitized string
	Make string
	Model string
}

func getRecords(filename string) ExifResult {
	var st C.struct_tagExifResult
	filename_c_str := C.CString(filename)
	C.print_exif(filename_c_str, &st)
	C.free(unsafe.Pointer(filename_c_str))
	return ExifResult {
		DateTime:           C.GoString(&st.DateTime[0]),
		DateTimeOriginal:   C.GoString(&st.DateTimeOriginal[0]),
		DateTimeDigitized:  C.GoString(&st.DateTimeDigitized[0]),
		Make:               C.GoString(&st.Make[0]),
		Model:              C.GoString(&st.Model[0]),
	}
}
