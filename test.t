Simple:
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