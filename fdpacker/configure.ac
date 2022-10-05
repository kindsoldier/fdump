AC_INIT(fdpacker,0.0.1)
AM_INIT_AUTOMAKE
AC_PREFIX_DEFAULT(/usr/local)

PACKAGE=fdpacker

AC_CHECK_PROG(HAVE_GO, go, true, false, /bin /usr/local/bin /usr/bin)
if test "x$HAVE_GO" = "xfalse"; then
    AC_MSG_ERROR([Requested program go not found])
fi
AC_PATH_PROG([go],[go])
AC_PATH_PROG([GO],[go])

AC_PATH_PROG([CP],[gcp cp])
if test -z "$CP"; then
    AC_MSG_ERROR([Requested program cp not found])
fi
AC_PATH_PROG([CHOWN],[gchown chown])
if test -z "$CHOWN"; then
    AC_MSG_ERROR([Requested program chown not found])
fi


AC_PROG_INSTALL
AC_CANONICAL_HOST

case $host_os in
    *freebsd* )
    AC_SUBST(ROOT_GROUP, "wheel")
    AM_CONDITIONAL(FREEBSD_OS, true)
    AM_CONDITIONAL(LINUX_OS, false)
    OSNAME=freebsd
    ROOT_GROUP=wheel
        ;;
    *linux* )
    AC_SUBST(ROOT_GROUP, "root")
    AM_CONDITIONAL(FREEBSD_OS, false)
    AM_CONDITIONAL(LINUX_OS, true)
    OSNAME=linux
    ROOT_GROUP=root
        ;;
esac

AM_CONDITIONAL(SYSTEMD, false)
if test -d /lib/systemd/system; then
    AM_CONDITIONAL(SYSTEMD, true)
fi

AC_ARG_ENABLE([devel-mode],
    AS_HELP_STRING([--enable-devel-mode], [Enable developmend mode]))

AC_DEFINE_UNQUOTED(srv_devel_mode, "false", [developmend mode])
AC_SUBST(srv_devel_mode, "false")

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    AC_DEFINE_UNQUOTED(srv_devel_mode, "true", [developmend mode])
    AC_SUBST(srv_devel_mode, "true")
    SRCDIR=`pwd`
    enable_devel_mode=yes
])


test "x$prefix" == "xNONE" && prefix=$ac_default_prefix
PREFIX=$prefix

dnl --------------------------------------------------------------------------------------

case $host_os in
    *freefdsd* )
    default_user="daemon"
    default_group="daemon"
        ;;
    *linux* )
    default_user="www-data"
    default_group="www-data"
        ;;
esac

AC_ARG_WITH(user,
    AS_HELP_STRING([--with-user=${PACKAGE}],[set executing user name]),
    [ if test ! -z "$with_user" ; then
        case $with_user in
            "")
                AC_MSG_ERROR(You must specify user name)
                ;;
            *)
                SRV_USER="$with_user"
                ;;
        esac
    else
        SRV_USER="$default_user"
    fi ],
    [ SRV_USER="$default_user" ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRV_USER="`id -un`"
])

AC_MSG_NOTICE(srv_user set as ${SRV_USER})

AC_DEFINE_UNQUOTED(srv_user, "$SRV_USER", [effective user])
AC_SUBST(srv_user, "$SRV_USER")

AC_DEFINE_UNQUOTED(SRV_USER, "$SRV_USER", [effective user])
AC_SUBST(SRV_USER, "$SRV_USER")


dnl --------------------------------------------------------------------------------------

AC_ARG_WITH(group,
    AS_HELP_STRING([--with-group=${PACKAGE}],[set executing group name]),
    [ if test ! -z "$with_group" ; then
        case $with_group in
            "")
                AC_MSG_ERROR(You must specify group name)
                ;;
            *)
                SRV_GROUP="$with_group"
                ;;
        esac
    else
        SRV_GROUP="$default_group"
    fi ],
    [ SRV_GROUP="$default_group" ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRV_GROUP="`id -gn`"
])

AC_MSG_NOTICE(srv_group set as ${SRV_GROUP})

AC_DEFINE_UNQUOTED(srv_group, "$SRV_GROUP", [effective group])
AC_SUBST(srv_group, "$SRV_GROUP")

AC_DEFINE_UNQUOTED(SRV_GROUP, "$SRV_GROUP", [effective group])
AC_SUBST(SRV_GROUP, "$SRV_GROUP")

dnl --------------------------------------------------------------------------------------


AC_OUTPUT([
Makefile
])
