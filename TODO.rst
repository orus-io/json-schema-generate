TODO
====

Potential improvements identified during codebase analysis.
Not prioritized for immediate work — kept as a reference.


Critical (bugs & panics)
------------------------

1. ``panic`` in public API — ``Output()`` panics on template errors
   (``output.go:117-127``), ``Schema.Init()`` panics on malformed schemas
   (``jsonschema.go:232-234``). Both should return ``error``.

2. Index-out-of-range panics on empty strings — ``generator.go:207``,
   ``output_tmpl.go:32,36,70`` access ``s[0:1]`` / ``t[0]`` without length
   checks.

3. Non-deterministic output — ``OutputData.OneOfs`` is a
   ``map[string]OneOf`` iterated non-deterministically in the template
   (``output.go:84``). Should be a sorted slice like ``Structs`` and
   ``Aliases``.

4. ``main.go`` exits 0 on output-file creation error —
   ``cmd/schema-generate/main.go:62-65`` uses ``return`` instead of
   ``os.Exit(1)``.

5. Output file never closed — ``main.go:60`` creates the file but never
   calls ``Close()`` or ``Sync()``.

6. ``log.Fatal`` in tests — ``test/additionalProperties2_test.go:84,99``
   kills the entire test binary instead of failing just the test.

7. Off-by-one in ``lineAndCharacter`` — ``input.go:56-78`` returns error
   when ``offset == len(bytes)``, which is a valid position for
   ``json.SyntaxError.Offset``.

8. ``mise.toml`` ``|| true`` swallows generation failures —
   ``mise.toml:28`` masks real errors in the generate task.

9. ``mise/config.toml`` cover task references ``./tests`` — directory is
   ``test/`` (singular), task would fail.

10. Inverted test assertion — ``test/additionalProperties2_test.go:83``
    calls ``log.Fatal`` when structs ARE deeply equal, with a message
    saying they are NOT.


High (correctness & error handling)
------------------------------------

11. Errors swallowed in ``generateOneOf`` — ``generator.go:197`` silently
    ignores ``GetSchemaByReference`` errors, may produce incorrect
    ``JSONType``.

12. Errors not wrapped with ``%w`` — ``input.go:18,23,39,47``,
    ``jsonschema.go:217`` use ``errors.New(err.Error())`` or string
    concatenation.

13. ``ParseWithSchemaKeyRequired`` returns partial schema on error —
    ``jsonschema.go:203`` returns ``s`` on parse error but ``nil`` on
    validation errors (lines 217, 220).

14. Misleading flag description — ``cmd/schema-generate/main.go:20`` says
    "Allow input files with no $schema key" but setting it to true
    REQUIRES the key.

15. Misleading error message — ``jsonschema.go:211`` says "must have a
    $schema key unless schemaKeyRequired flag is set" — backwards.

16. ``updatePathElements`` and ``updateURIs`` skip
    ``AnyOf``/``AllOf``/``OneOf`` — ``jsonschema.go:238-262``,
    ``refresolver.go:100-165`` don't traverse these fields, but
    ``updateParentLinks`` does. References inside them can't be resolved.

17. ``AnyOf``/``AllOf`` silently ignored by generator —
    ``generator.go:101-163`` only handles ``OneOf``, not
    ``AnyOf``/``AllOf``.

18. ``getPrimitiveTypeName`` returns error-indicating strings —
    ``generator.go:378,391,401`` returns ``"error_creating_array"`` etc.
    as type names alongside errors.


Medium (code quality & maintainability)
---------------------------------------

19. Triplicated generation logic — ``test/gen.sh``, ``test/gen.go``, and
    ``mise.toml`` all implement the same loop. Consolidate to one source
    of truth.

20. Monolithic 700-line template — ``output_tmpl.go:95-811`` should be
    split into sub-templates.

21. ``processSchema`` overly complex — ``generator.go:101-163`` handles
    too many concerns. Should be decomposed.

22. ``processObject`` does too many things — ``generator.go:277-368``,
    ~90 lines with multiple responsibilities.

23. ``OneOf*Null`` types are copy-pasted — ``output_tmpl.go:220-464``
    three identical types differing only in value type. Could use
    generics.

24. ``getOrderedFieldNames`` and ``getOrderedStructNames`` are identical
    — ``output.go:10-30``, could use ``slices.Sorted(maps.Keys(m))``.

25. ``Type()`` and ``MultiType()`` share parsing logic —
    ``jsonschema.go:141-182``, could extract a common helper.

26. ``updateParentLinks`` and ``updatePathElements`` have parallel
    structure — ``jsonschema.go:238-300``, risk of one being updated but
    not the other (already happened with ``AnyOf``/``AllOf``).

27. ``AdditionalProperties`` type alias requires frequent casts —
    ``jsonschema.go:14``, cast everywhere with ``(*Schema)(...)``.

28. ``contains`` wrapper is unnecessary — ``generator.go:370-372`` wraps
    ``slices.Contains`` trivially.

29. Stale ``go.sum`` entry — ``go.sum:3-4`` has orphaned ``go-orusapi``
    module. Run ``go mod tidy``.

30. Outdated dependency — ``go.mod:5`` ``shopspring/decimal`` is a 2019
    pseudo-version, should use tagged release.

31. ``Generator`` exposes mutable internal state — ``Structs``,
    ``Aliases``, ``OneOfs`` are exported map fields
    (``generator.go:19-21``).


Low (naming, docs, cleanup)
---------------------------

32. Misspelled identifiers — ``output_tmpl.go:16``
    ``iteratorUnmashallerTypes``, ``output_tmpl.go:39`` ``deferedType``.

33. Wrong comment text — ``output_tmpl.go:322`` says "string" should say
    "number", ``output_tmpl.go:405`` says "number" should say "bool".

34. Single-letter flag variables — ``cmd/schema-generate/main.go:15-17``
    (``o``, ``p``, ``i``).

35. Dead code: ``NameCount`` — ``jsonschema.go:75`` never used.

36. Dead code: ``Examples`` — ``jsonschema.go:63`` never used.

37. Dead code: ``GetByJSONType`` — ``generator.go:519`` never called.

38. Dead code: ``EmptyNumber`` in marshaller lists —
    ``output_tmpl.go:10,17``, never generated.

39. Empty test functions —
    ``test/additionalPropertiesMarshal_test.go:10-17`` three tests with
    empty bodies.

40. Commented-out debug code — 10 ``//Output(...)`` lines in
    ``generator_test.go``.

41. Undocumented exported fields — ``Struct.Fields``,
    ``Struct.GenerateCode``, ``Struct.AdditionalType``, ``Field.Enum``,
    ``OneOfType.*``, ``OutputData.*``.

42. ``vega-lite-v2.0.json_`` — 270KB tracked file with trailing
    underscore to skip glob, undocumented convention.

43. Empty ``build/`` directory — untracked leftover artifact.

44. Deprecated ``-i`` flag without deprecation notice —
    ``cmd/schema-generate/main.go:17``.

45. Missing ``--version`` flag on CLI.

46. No ``vulncheck`` in CI — configured in mise but not in
    ``.github/workflows/ci.yml``.


Test coverage gaps
------------------

47. ``cmd/schema-generate/main.go`` — zero tests (main_test.go is
    empty).

48. ``Output()`` never directly tested — only via integration.

49. ``ReadInputFiles`` error paths untested — file errors, JSON syntax
    errors, unmarshal type errors.

50. ``RefResolver`` methods untested — ``GetPath``,
    ``GetSchemaByReference``, ``InsertURI``.

51. ``anyOf``/``allOf`` completely untested — no schemas, no tests.

52. 17 of 24 generated packages have no behavioral tests.

53. ``recursion.json`` — circular references generated but never
    round-trip tested.

54. ``multipleOf`` with decimal — branch in ``processSchema`` never
    tested.

55. Tests that could be table-driven — ``TestFieldGeneration``,
    ``TestNestedStructGeneration``/``TestEmptyNestedStructGeneration``
    (near-identical), ``TestAp*`` tests.

56. Fragile count assertions — ``generator_test.go:84,160,210,251,355,
    413,495`` hardcode exact struct counts.

57. ``marshal_test.go`` uses hand-written types — not generated code,
    doesn't validate actual output.

58. ``TestThatUnmarshallingIsPossible`` tests wrong package —
    ``generator_test.go:628`` uses hand-written structs, not generated
    ones.
