defmodule Variables do
  def demo do
    # ===== ПЕРЕМЕННЫЕ =====
    x = 42
    IO.puts("x = #{x}")

    # Переприсваивание (создание новой переменной)
    x = x + 1
    IO.puts("x после переприсваивания = #{x}")

    # ===== ЧИСЛА =====
    integer = 255
    hex = 0xFF
    binary = 0b11111111
    octal = 0o377
    float = 3.14
    scientific = 1.0e-10

    IO.puts("\n--- Числа ---")
    IO.puts("Integer: #{integer}")
    IO.puts("Hex: #{hex}")
    IO.puts("Binary: #{binary}")
    IO.puts("Octal: #{octal}")
    IO.puts("Float: #{float}")
    IO.puts("Scientific: #{scientific}")

    # ===== АТОМЫ =====
    atom1 = :hello
    atom2 = :ok
    atom3 = :error
    # это :true
    bool_true = true
    # это :false
    bool_false = false
    # это :nil
    null = nil

    IO.puts("\n--- Атомы ---")
    IO.puts("Atom: #{atom1}")
    IO.puts("OK: #{atom2}")
    IO.puts("Error: #{atom3}")
    IO.puts("True: #{bool_true}")
    IO.puts("False: #{bool_false}")
    IO.puts("Nil: #{inspect(null)}")

    # ===== СТРОКИ =====
    string1 = "Hello, World!"
    string2 = "Привет, мир!"
    name = "Elixir"
    interpolated = "Hello, #{name}!"

    multiline = """
    Это
    многострочная
    строка
    """

    IO.puts("\n--- Строки ---")
    IO.puts(string1)
    IO.puts(string2)
    IO.puts(interpolated)
    IO.puts("Длина строки '#{name}': #{String.length(name)}")

    # ===== СПИСКИ =====
    list1 = [1, 2, 3, 4, 5]
    list2 = [1, "two", :three, 4.0]
    # добавление в начало
    list3 = [0 | list1]

    IO.puts("\n--- Списки ---")
    IO.inspect(list1, label: "List1")
    IO.inspect(list2, label: "Mixed list")
    IO.inspect(list3, label: "List with cons")

    # Pattern matching со списками
    [head | tail] = list1
    IO.puts("Head: #{head}")
    IO.inspect(tail, label: "Tail")

    # ===== КОРТЕЖИ =====
    tuple1 = {1, 2, 3}
    tuple2 = {:ok, "success"}
    tuple3 = {:error, "failed", 404}

    IO.puts("\n--- Кортежи ---")
    IO.inspect(tuple1, label: "Tuple1")
    IO.inspect(tuple2, label: "OK tuple")
    IO.inspect(tuple3, label: "Error tuple")
    IO.puts("Элемент по индексу 1: #{elem(tuple1, 1)}")
    IO.puts("Размер кортежа: #{tuple_size(tuple1)}")

    # ===== КАРТЫ (MAPS) =====
    map1 = %{:a => 1, :b => 2, :c => 3}
    map2 = %{name: "John", age: 30, city: "Moscow"}

    IO.puts("\n--- Карты ---")
    IO.inspect(map1, label: "Map1")
    IO.inspect(map2, label: "Map2")
    IO.puts("Доступ к ключу :name = #{map2.name}")
    IO.puts("Доступ через [] = #{map2[:age]}")

    # Обновление карты
    map3 = %{map2 | age: 31}
    IO.inspect(map3, label: "Updated map")

    # ===== KEYWORD LISTS =====
    keyword = [name: "Alice", age: 25, city: "SPb"]

    IO.puts("\n--- Keyword Lists ---")
    IO.inspect(keyword, label: "Keyword list")
    IO.puts("Значение :name = #{keyword[:name]}")

    # ===== ПРОВЕРКА ТИПОВ =====
    IO.puts("\n--- Проверка типов ---")
    IO.puts("is_integer(42): #{is_integer(42)}")
    IO.puts("is_float(3.14): #{is_float(3.14)}")
    IO.puts("is_number(42): #{is_number(42)}")
    IO.puts("is_atom(:hello): #{is_atom(:hello)}")
    IO.puts("is_boolean(true): #{is_boolean(true)}")
    IO.puts("is_binary(\"hello\"): #{is_binary("hello")}")
    IO.puts("is_list([1,2,3]): #{is_list([1, 2, 3])}")
    IO.puts("is_tuple({1,2}): #{is_tuple({1, 2})}")
    IO.puts("is_map(%{a: 1}): #{is_map(%{a: 1})}")

    # ===== АНОНИМНЫЕ ФУНКЦИИ =====
    IO.puts("\n--- Функции ---")
    add = fn a, b -> a + b end
    IO.puts("Анонимная функция add.(5, 3) = #{add.(5, 3)}")

    # Короткий синтаксис
    multiply = &(&1 * &2)
    IO.puts("Короткий синтаксис multiply.(4, 5) = #{multiply.(4, 5)}")

    # ===== PATTERN MATCHING =====
    IO.puts("\n--- Pattern Matching ---")
    {status, message} = {:ok, "Success"}
    IO.puts("Status: #{status}, Message: #{message}")

    [first, second | rest] = [1, 2, 3, 4, 5]
    IO.puts("First: #{first}, Second: #{second}")
    IO.inspect(rest, label: "Rest")

    # ===== ДИАПАЗОНЫ =====
    range = 1..10
    IO.puts("\n--- Диапазоны ---")
    IO.inspect(range, label: "Range")
    IO.puts("3 в диапазоне? #{3 in range}")

    # ===== ОПЕРАТОРЫ =====
    IO.puts("\n--- Операторы ---")
    IO.puts("5 + 3 = #{5 + 3}")
    IO.puts("10 - 4 = #{10 - 4}")
    IO.puts("6 * 7 = #{6 * 7}")
    IO.puts("20 / 5 = #{20 / 5}")
    IO.puts("div(17, 5) = #{div(17, 5)}")
    IO.puts("rem(17, 5) = #{rem(17, 5)}")

    IO.puts("\"hello\" <> \" world\" = # {\"hello\" <> \" world\"}")
    IO.puts("[1,2] ++ [3,4] = #{inspect([1, 2] ++ [3, 4])}")

    IO.puts("\n--- Логические операторы ---")
    IO.puts("true and false = #{true and false}")
    IO.puts("true or false = #{true or false}")
    IO.puts("not true = #{not true}")

    IO.puts("\n--- Операторы сравнения ---")
    IO.puts("5 == 5.0 = #{5 == 5.0}")
    IO.puts("5 === 5.0 = #{5 === 5.0}")
    IO.puts("5 != 6 = #{5 != 6}")
    IO.puts("5 < 10 = #{5 < 10}")

    :ok
  end

  # Именованные функции с @spec
  @spec add(integer, integer) :: integer
  def add(a, b), do: a + b

  @spec multiply(number, number) :: number
  def multiply(a, b), do: a * b

  @spec greet(String.t()) :: String.t()
  def greet(name) do
    "Hello, #{name}!"
  end
end

Variables.demo()

IO.puts("\n--- Вызов именованных функций ---")
IO.puts("add(10, 20) = #{Variables.add(10, 20)}")
IO.puts("multiply(3, 7) = #{Variables.multiply(3, 7)}")
IO.puts(Variables.greet("Golang Developer"))
