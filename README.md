# Togo4
    Simple console app for managing my todos, in order to make me go for them.
# Language: GO

# Hint
# +: New Togo:
... +   title   [=  description]    [+p   progress_till_now]   [+w   weight_of_togo]    [+x | -x]   [@  start_date_as_how_many_days_from_now    start_time_as_hh:mm]    [NEXT_COMMAND]
# #: Show Togos
... # [NEXT_COMMAND]
# %: Progress Made:
... % [NEXT_COMMAND]

    Calculate the progress been made till now:

*   Tags order is optional, and tags and their param must be seperated by tabs.
*   Each line can contain multiple command, as many as you want. for 4x:
        +   new_togo    @   1   10:00   +p  85  #  +   next_togo   +x  #   %
*   Extra:
        +x: its an extra Togo. its not mandatory but has extra points doing it.
        -x: not extra (default)
*   all params between [] are optional.

# How to use:
    first install Togo package
    then:

    go mod init togo4
    go mod tidy
    go run togo4.go
