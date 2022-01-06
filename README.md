# programowanie współbieżne

każda lista zawiera rozwiązanie zadania w języku go i ada.
* listy 1 nie wrzucałem - druga lista ją rozszeża
* nie ma listy 4 w języku ada

## lista 1

Program ma być uruchamiany z parametrami: n, d, b

generuje graf oraz pakiety i przesyła współbieżnie pakiety do docelowego wierzchołka

parametry:
* n - wielkość grafu (graf zawiera tylko ścieżkę hamiltona)
* d - liczba dodatkowych krawędzi w grafie
* b - liczba pakietów

## lista 2

Program ma być uruchamiany z parametrami: n, d, k, b, h, klusownik

rozszeża listę 1

paramety:
* n - wielkość grafu (graf zawiera tylko ścieżkę hamiltona)
* d - liczba dodatkowych krawędzi w grafie
* b - liczba pakietów
* k - maksymalna liczba skoków dla każdego parametru
* klusownik - po losowym czasie losowo ustawia sie we wierzchołku i czeka na pakiet po czym go usuwa

## lista 3

protokół routingu podobny do znanego protokołu RIP

parametry:
* n - wielkość grafu (graf zawiera tylko ścieżkę hamiltona)
* d - liczba dodatkowych krawędzi w grafie

## lista 4

rozszeżenie listy 3 o dodanie hostów, algorytm podczas wysyłania pakietów wybiera najkrótsze ścieżki między hostami dzięki czemu im więcej pakietów zostaje wysłanych przemieszczają się one w coraz mniejszej ilości skoków

parametry:
* n - wielkość grafu (graf zawiera tylko ścieżkę hamiltona)
* d - liczba dodatkowych krawędzi w grafie
* h - liczba hostów
* 
