AC_INIT(fdagent,0.0.1)
AM_INIT_AUTOMAKE
AC_PREFIX_DEFAULT(/usr/local)

PACKAGE=fdagentd

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

AC_ARG_WITH(confdir,
    AS_HELP_STRING([--with-confdir=PATH],[set configuration dir to PATH (default: "${ac_default_prefix}"/etc/${PACKAGE})]),
    [ if test ! -z "$with_confdir" ; then
        case $with_confdir in
            /*)
                SRV_CONFDIR="$with_confdir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an absolute path to --with-confdir=PATH)
                ;;
        esac
    else
        SRV_CONFDIR="$prefix/etc/${PACKAGE}"
    fi ],
    [
    SRV_CONFDIR="$prefix/etc/${PACKAGE}"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRV_CONFDIR="${SRCDIR}/"
    sysconfdir="${SRCDIR}/"
], [
    test "x$SRV_CONFDIR" == "x/usr/etc/${PACKAGE}" && SRV_CONFDIR="/etc/${PACKAGE}"
    test "x$prefix" == "x/usr" && sysconfdir="/etc"
])

AC_MSG_NOTICE(srv_confdir set as ${SRV_CONFDIR})

dnl AC_DEFINE_UNQUOTED(SRV_CONFDIR, "$SRV_CONFDIR", [location of configuration files for ${PACKAGE}])
dnl AC_DEFINE_UNQUOTED(srv_confdir, "$SRV_CONFDIR", [location of configuration files for ${PACKAGE}])
AC_SUBST(SRV_CONFDIR, "$SRV_CONFDIR")
AC_SUBST(srv_confdir, "$SRV_CONFDIR")

dnl --------------------------------------------------------------------------------------

AC_ARG_WITH(logdir,
    AS_HELP_STRING([--with-logdir=PATH],[set file path for source logdir (default: /var/log/${PACKAGE}/${PACKAGE}.log)]),
    [ if test ! -z "$with_logdir" ; then
        case $with_logdir in
            /*)
                SRV_LOGDIR="$with_logdir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an afdaolute path to --with-logdir=PATH)
                ;;
        esac
    else
        SRV_LOGDIR="/var/log/${PACKAGE}"
    fi ],
    [
    SRV_LOGDIR="/var/log/${PACKAGE}"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRV_LOGDIR="${SRCDIR}/tmp.log"
])

AC_MSG_NOTICE(srv_logdir set as ${SRV_LOGDIR})

dnl AC_DEFINE_UNQUOTED(SRV_LOGDIR, "$SRV_LOGDIR", [location of ${PACKAGE} logdir])
dnl AC_DEFINE_UNQUOTED(srv_logdir, "$SRV_LOGDIR", [location of ${PACKAGE} logdir])
AC_SUBST(SRV_LOGDIR, "$SRV_LOGDIR")
AC_SUBST(srv_logdir, "$SRV_LOGDIR")

dnl --------------------------------------------------------------------------------------

AC_ARG_WITH(rundir,
    AS_HELP_STRING([--with-rundir=PATH],[set file path for source rundir (default: /var/run/${PACKAGE})]),
    [ if test ! -z "$with_rundir" ; then
        case $with_rundir in
            /*)
                SRV_RUNDIR="$with_rundir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an afdaolute path to --with-rundir=PATH)
                ;;
        esac
    else
        SRV_RUNDIR="/var/run/${PACKAGE}"
    fi ],
    [
    SRV_RUNDIR="/var/run/${PACKAGE}"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRV_RUNDIR="${SRCDIR}/tmp.run"
])

AC_MSG_NOTICE(srv_rundir set as ${SRV_RUNDIR})

dnl AC_DEFINE_UNQUOTED(SRV_RUNDIR, "$SRV_RUNDIR", [location of pid file])
dnl AC_DEFINE_UNQUOTED(srv_rundir, "$SRV_RUNDIR", [location of pid file])
AC_SUBST(SRV_RUNDIR, "$SRV_RUNDIR")
AC_SUBST(srv_rundir, "$SRV_RUNDIR")

dnl --------------------------------------------------------------------------------------

default_datadir="/var/db/$PACKAGE"

AC_ARG_WITH(datadir,
    AS_HELP_STRING([--with-datadir=PATH],[set data directory (default: $default_datadir)]),
    [ if test ! -z "$with_datadir" ; then
        case $with_datadir in
            /*)
                SRV_DATADIR="$with_datadir"
                ;;
            *)
                                AC_MSG_ERROR(You must specify an afdaolute path to --with-datadir=PATH)
                ;;
        esac
    else
        SRV_DATADIR="$default_datadir"
    fi ],
    [
    SRV_DATADIR="$default_datadir"
    ])

AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRV_DATADIR="${SRCDIR}/tmp.data"
])

AC_MSG_NOTICE(srv_datadir set as ${SRV_DATADIR})

dnl AC_DEFINE_UNQUOTED(SRV_DATADIR, "$SRV_DATADIR", [location of data dir])
dnl AC_DEFINE_UNQUOTED(srv_datadir, "$SRV_DATADIR", [location of data dir])
AC_SUBST(SRV_DATADIR, "$SRV_DATADIR")
AC_SUBST(srv_datadir, "$SRV_DATADIR")



dnl --------------------------------------------------------------------------------------

case $host_os in
    *freefdad* )
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


SRV_BACKGROUND="true"
AS_IF([test "x$enable_devel_mode" = "xyes"], [
    SRV_BACKGROUND="false"
])

AC_SUBST(srv_background, "$SRV_BACKGROUND")

AC_SUBST(srv_conf, "$PACKAGE.conf")
AC_SUBST(srv_name, "$PACKAGE")
AC_SUBST(srv_port, 5300)

AC_OUTPUT([
Makefile
rc/Makefile
fdasrv/fdaconf.go
rc/rc.d/fdagentd
])
