AC_INIT(fdump,0.0.1)
AM_INIT_AUTOMAKE([foreign subdir-objects])

AC_ARG_ENABLE([fdstore],
    AS_HELP_STRING([--disable-fdstore], [disable fdstore module]),,
    AC_CONFIG_SUBDIRS([fdstore])
)

AC_ARG_ENABLE([fdagent],
    AS_HELP_STRING([--disable-fdagent], [disable fdagent module]),,
    AC_CONFIG_SUBDIRS([fdagent])
)

AC_ARG_ENABLE([fdmaster],
    AS_HELP_STRING([--disable-fdmaster], [disable fdmaster module]),,
    AC_CONFIG_SUBDIRS([fdmaster])
)


AC_CONFIG_FILES(Makefile)
AC_OUTPUT
