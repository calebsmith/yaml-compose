Noop - Output for a normal YAML file should be invariant:
  $ pushd $INPUTS > /dev/null; yaml-compose vars.yaml; popd > /dev/null
  A: 1
  B: 2
  data:
    - C: 3
    - D: 4

Inject:
  $ pushd $INPUTS > /dev/null; yaml-compose inject.yaml; popd > /dev/null
  data:
    keya1: val1
    keya2: val2
    var3key: value

Load:
  $ pushd $INPUTS > /dev/null; yaml-compose load.yaml; popd > /dev/null
  Value: 1

Full:
  $ pushd $INPUTS > /dev/null; yaml-compose main.yaml; popd > /dev/null
  array_of_maps:
    - keya1: vala1
      keya2: vala2
    - keyb1: valb1
      keyb2: valb2
    - inline_map:
        keya1: val1
        keya2: val2
  rendered:
    - 1
    - 2
    - value
    - value1
    - value2
    - value3
    - new-value
    - maps:
      C: 3
      D: 4