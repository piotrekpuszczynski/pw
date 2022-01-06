--Piotr Puszczynski 

with
    Ada.Text_IO,
    Ada.Command_Line,
    Ada.Strings.Bounded,
    Ada.numerics.discrete_random,
    GNAT.OS_Lib;
use 
    Ada.Text_IO,
    Ada.Command_Line;


procedure lista3 is
	
	n : Integer := 5;
	d : Integer := 2;

	type randRange is new Integer range 0 .. (n - 1);
	package Rand_Int is new ada.numerics.discrete_random(randRange);
   	use Rand_Int;
   	gen : Generator;
	
	procedure checkData is
	begin
		if n < 2 or d < 0 or d > ((n - 3) * n / 2) + 1 then
			Put_Line("invalid arguments");
			GNAT.OS_Lib.OS_Exit (0);
		end if;
	end checkData;
--------------------
	type lengthArray is array (0 .. (n - 1)) of Integer;
	type boolLengthArray is array (0 .. (n - 1)) of Boolean;

	type Packet is record
		j : lengthArray;
		cost : lengthArray;
		v : Integer;
	end record;

	type Vertex is record
		id : Integer;
		connections : lengthArray;
		p : Packet;
		cost : lengthArray;
		nextHop : lengthArray;
		changed : boolLengthArray;
	end record;
--------------------
	vertexes : array (0 .. (n - 1)) of Vertex;

	procedure initializeVertexes is
	i : Integer := 0;
	emptyInt : lengthArray;
	emptyP : Packet;
	emptyBool : boolLengthArray;
	begin
		for j in emptyInt'Range loop
			emptyInt(j) := -1;
			emptyBool(j) := false;
		end loop;
		emptyP := (j => emptyInt, cost => emptyInt, v => -1);

		while i < n loop
			vertexes(i) := (id => i, connections => emptyInt, p => emptyP, cost => emptyInt, nextHop => emptyInt, changed => emptyBool);
			i := i + 1;
		end loop;
	end initializeVertexes;

	procedure initializeShortcuts is
	i : Integer := 0;
	rand1 : Integer;
	rand2 : Integer;
	correctRandoms : Boolean;
	begin
		while i < n - 1 loop
			for j in vertexes(i).cost'Range loop
				if vertexes(i).connections(j) = -1 then
					vertexes(i).connections(j) := i + 1;
					--Put_Line(Integer'Image(i) & " - " & Integer'Image(i + 1));
					exit when true;
				end if;
			end loop;
			i := i + 1;
		end loop;

		i := 1;
		while i < n loop
			for j in vertexes(i).cost'Range loop
				if vertexes(i).connections(j) = -1 then
					vertexes(i).connections(j) := i - 1;
					--Put_Line(Integer'Image(i) & " - " & Integer'Image(i - 1));
					exit when true;
				end if;
			end loop;
			i := i + 1;
		end loop;

		i := 0;
		while i < d loop
			reset(gen);
			rand1 := Integer(random(gen));
			reset(gen);
			rand2 := Integer(random(gen));
			correctRandoms := false;
			
			if rand1 /= rand2 then
				for j in vertexes(rand1).connections'Range loop
					if vertexes(rand1).connections(j) = rand2 then
						correctRandoms := true;
						exit when true;
					end if;
				end loop;
			else
				correctRandoms := true;
			end if;

			if not correctRandoms then
				for j in vertexes(rand1).connections'Range loop
					if vertexes(rand1).connections(j) = -1 then
						vertexes(rand1).connections(j) := rand2;
						--Put_Line(Integer'Image(rand1) & " -" & Integer'Image(rand2));
						exit when true;
					end if;
				end loop;

				for j in vertexes(rand2).connections'Range loop
					if vertexes(rand2).connections(j) = -1 then
						vertexes(rand2).connections(j) := rand1;
						--Put_Line(Integer'Image(rand2) & " -" & Integer'Image(rand1));
						exit when true;
					end if;
				end loop;
			else
				i := i - 1;
			end if;

			i := i + 1;
		end loop;

	end initializeShortcuts;

	procedure initializeRoutingTable is
	i : Integer := 0;
	j : Integer := 0;
	k : Integer;
	contains : Boolean;
	begin
		while i < n loop
			j := 0;
			while j < n loop
				if i /= j then
					contains := false;

					k := 0;
					while k < n loop
						if vertexes(i).connections(k) /= -1 then
							if vertexes(vertexes(i).connections(k)).id = j then
								contains := true;
								exit when true;
							end if;
						end if;
						k := k + 1;
					end loop;

					if contains then
						vertexes(i).cost(j) := 1;
						vertexes(i).nextHop(j) := j;
					else
						if i - j < 0 then
							vertexes(i).cost(j) := j - i;
							vertexes(i).nextHop(j) := i - 1;
						else
							vertexes(i).cost(j) := i - j;
							vertexes(i).nextHop(j) := i + 1;
						end if;

					end if;
					vertexes(i).changed(j) := true;
				end if;
				j := j + 1;
			end loop;
			i := i + 1;
		end loop;

	end initializeRoutingTable;

	procedure drawGraph is

		procedure draw(i : Integer; j : Integer) is
			a : Integer := 0;
		begin
			while a < i loop
				Put("      ");
				a := a + 1; 
			end loop;
			Put(Integer'Image(i) & " ");
			while a < j - 1 loop
				Put("----");
				a := a + 1;
			end loop;
			Put("---" & Integer'Image(j));
			New_Line;
		end draw;

	i : Integer := 0;
	j : Integer;
	was : Boolean;
	begin
		while i < n loop
			j := 0;
			while j < n loop
				was := false;
				for k in vertexes(i).connections'Range loop
					if vertexes(i).connections(k) = j then
						was := true;
						exit when true;
					end if;
				end loop;

				if i < j and was then
					draw(i, j);
				end if;
				j := j + 1;
			end loop;
			i := i + 1;
		end loop;
	end drawGraph;
--------------------
	task type sender is
		entry Start (j : Integer);
	end sender;

	task body sender is
	taskId : Integer;
	r : Float;
	jTab : lengthArray;
	costTab : lengthArray;
	send : Boolean;
	begin
		accept Start (j : Integer) do
			taskId := j;
		end Start;

		while true loop
			reset(gen);
			r := Float(random(gen)) / Float(n);
			delay Duration(r);

			for j in jTab'Range loop
				jTab(j) := -1;
				costTab(j) := -1;
			end loop;

			send := false;
			for c in vertexes(taskId).changed'Range loop
				if vertexes(taskId).changed(c) then
					send := true;
					--Put_Line(Integer'Image(taskId));
					exit when true;
				end if;
			end loop;

			for i in vertexes(taskId).cost'Range loop
				if vertexes(taskId).changed(i) then
					for k in jTab'Range loop
						if jTab(k) = -1 then
							jTab(k) := i;
							costTab(k) := vertexes(taskId).cost(i);
							vertexes(taskId).changed(i) := false;
						end if;
					end loop;
				end if;
			end loop;

			if send then
				for j in vertexes(taskId).connections'Range loop
					if vertexes(taskId).connections(j) /= -1 then
						while vertexes(vertexes(taskId).connections(j)).p.v /= -1 loop
							reset(gen);
							r := 0.5 * Float(random(gen)) / Float(n);
							delay Duration(r);
						end loop;
						vertexes(vertexes(taskId).connections(j)).p := (j => jTab, cost => costTab, v => taskId);
						Put_Line("vertex" & Integer'Image(taskId) & " sent to" & Integer'Image(vertexes(taskId).connections(j)));
					else
						exit when true;
					end if;
				end loop;
			end if;
		end loop;
	end sender;
--------------------
	task type receiver is
		entry Start (j : Integer);
	end receiver;

	task body receiver is
	taskId : Integer;
	r : Float;
	emptyInt : lengthArray;
	newCost : Integer;
	begin
		accept Start (j : Integer) do
			taskId := j;
		end Start;
		while true loop
			while vertexes(taskId).p.v = -1 loop
				reset(gen);
				r := 0.5 * Float(random(gen)) / Float(n);
				delay Duration(r);
			end loop;

			Put_Line("vertex" & Integer'Image(taskId) & " received from" & Integer'Image(vertexes(taskId).p.v));

			for j in emptyInt'Range loop
				exit when vertexes(taskId).p.v = -1;
				
				newCost := 1 + vertexes(taskId).p.cost(j);

				if newCost < vertexes(vertexes(taskId).p.j(j)).cost(taskId) then
					vertexes(vertexes(taskId).p.j(j)).cost(taskId) := newCost;
					vertexes(vertexes(taskId).p.j(j)).nextHop(taskId) := vertexes(vertexes(taskId).p.v).id;
					vertexes(vertexes(taskId).p.j(j)).changed(taskId) := true;
				end if;
			end loop;

			for j in emptyInt'Range loop
				emptyInt(j) := -1;
			end loop;

			vertexes(taskId).p := (j => emptyInt, cost => emptyInt, v => -1);
		end loop;

	end receiver;

	task stop;
	task body stop is
	t : Boolean := true;
	temp : Boolean;
	begin
		while t loop
			delay 1.0;
			temp := false;
			for j in vertexes'Range loop
				for i in vertexes(j).cost'Range loop
					if vertexes(j).changed(i) then
						temp := true;
					end if;
				end loop;
			end loop;
			t := temp;
		end loop;

		for j in vertexes'Range loop
			for i in vertexes(j).cost'Range loop
				if vertexes(j).cost(i) = -1 then
					Put(Integer'Image(0) & " ");
				else
					Put(Integer'Image(vertexes(j).cost(i)) & " ");
				end if;

			end loop;
			New_Line(1);
		end loop;

		GNAT.OS_Lib.OS_Exit (0);
	end stop;

	
	senders : array (0 .. (n - 1)) of sender;
	receivers : array (0 .. (n - 1)) of receiver;

	procedure start is
	begin
		checkData;
		initializeVertexes;
		initializeShortcuts;
		initializeRoutingTable;
		drawGraph;

		for i in senders'Range loop
			senders(i).Start(i);
		end loop;
		for i in receivers'Range loop
			receivers(i).Start(i);
		end loop;

	end start;

begin
	start;
end lista3;
