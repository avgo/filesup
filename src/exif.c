#include <stdio.h>

#include <libexif/exif-loader.h>

#define TAG_VALUE_BUF    1024

static int print_exif(const char *filename);
void print_usage(const char *cmd);
static void show1(ExifContent *c, void *data);
static void show2(ExifEntry *e, void *data);

static int print_exif(const char *filename)
{
	printf("file: %s\n", filename);

	ExifLoader *l = exif_loader_new();

	exif_loader_write_file(l, filename);

	ExifData *ed = exif_loader_get_data(l);

	if (ed == NULL)
	{
		fprintf(stderr, "error: some error occurs\n");
		return -1;
	}

	exif_data_foreach_content(ed, show1, NULL);
}

static void show1(ExifContent *c, void *data)
{
	exif_content_foreach_entry(c, show2, data);
}

static void show2(ExifEntry *e, void *data)
{
	ExifIfd ifd = exif_entry_get_ifd(e);
	const char *str = exif_tag_get_title_in_ifd(e->tag, ifd);
	if (!(e->tag == EXIF_TAG_DATE_TIME ||
		e->tag == EXIF_TAG_DATE_TIME_ORIGINAL ||
		e->tag == EXIF_TAG_DATE_TIME_DIGITIZED ||
		e->tag == EXIF_TAG_MAKE ||
		e->tag == EXIF_TAG_MODEL))
		return;
	char val_buf[TAG_VALUE_BUF];
	const char *val = exif_entry_get_value(e, val_buf, sizeof val_buf);
	printf("%s: '%s'\n", str, val);
}

