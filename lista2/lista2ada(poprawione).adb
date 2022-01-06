--Piotr Puszczynski
--parametry mozna wprowadzic jedynie przed kompilacja
with
    Ada.Text_IO,
    Ada.Command_Line,
    Ada.Strings.Bounded,
    Ada.numerics.discrete_random,
    GNAT.OS_Lib;
use 
    Ada.Text_IO,
    Ada.Command_Line;


procedure Main is
	
	n : Integer := 5;
	d : Integer := 0;
	k : Integer := 3;
	b : Integer := 1;
	h : Integer := 5;
	kk : Boolean := true;
--------------------------------
	lastDone : Integer := -1;
	counter : Integer := k;

	type randRange is new Integer range 1 .. (n - 1);
	package Rand_Int is new ada.numerics.discrete_random(randRange);
   	use Rand_Int;
   	gen : Generator;
--------------------------------
	procedure checkD is
		i : Integer := 0;
		sum : Integer := 0;
	begin
		while n - 1 > i loop
			
			sum := sum + n - i - 1;
			i := i + 1;
			
		end loop;

		if d > sum then
			Put_Line("invalid arguments");
			GNAT.OS_Lib.OS_Exit (0);
		end if;
	end checkD;

	procedure checkB is
		i : Integer := 0;
		sum : Integer := 0;
	begin
		while n - 1 > i loop
			
			sum := sum + i;
			i := i + 1;
			
		end loop;

		if b > sum then
			Put_Line("invalid arguments");
			GNAT.OS_Lib.OS_Exit (0);
		end if;
	end checkB;
--------------------------------
	arrayOfRandom : array (0 .. (n - 1)) of Integer;

	procedure generateRandomArray is
		i : Integer;
		r : Integer;
	begin
		while d > 0 loop
			i := 0;
			while i < n - 2 loop
				reset(gen);
				r := Integer(random(gen)) / 4;

				if r + arrayOfRandom(i) > n - 2 - i then
					r := 0;
				end if;
				if r > d then
					r := d;
				end if;

				arrayOfRandom(i) := arrayOfRandom(i) + r;
				d := d - r;
				i := i + 1;
			end loop;
		end loop;
	end generateRandomArray;

	arrayOfReverseRandom : array (0 .. (n - 1)) of Integer;

	procedure generateReverseRandomArray is
		i : Integer;
		r : Integer;
		rand : Integer;
	begin
		while b > 0 loop
			i := 1;
			while i < n - 1 loop
				reset(gen);
				rand := Integer(random(gen));
				if rand > n / 2 then
					r := 1;
				else
					r := 0;
				end if;

				if r + arrayOfReverseRandom(i) > i then
					r := 0;
				end if;
				if r > b then
					r := 0;
				end if;

				arrayOfReverseRandom(i) := arrayOfReverseRandom(i) + r;
				b := b - r;
				i := i + 1;
			end loop;
		end loop;
	end generateReverseRandomArray;
--------------------------------
	type connectionsArray is array (0 .. (n - 2)) of Integer;
	type servicedPackets is array (0 .. (k * n)) of Integer;

	type Node is record
		id : Integer;
		packet : Integer;
		connections : connectionsArray;
		serviced : servicedPackets;
		trap: Boolean;
	end record;
--------------------------------
	type visitedNodes is array (0 .. (n * k)) of Integer;
	v : array (0 .. (k - 1)) of visitedNodes;
	packets : array (0 .. (k - 1)) of Integer;
	packetLifetime : array (0 .. (k - 1)) of Integer;

	procedure addPackets is
	begin
		for i in 0 .. (k - 1) loop
			packets(i) := i;
			packetLifetime(i) := 0;
		end loop;

		for i in v'Range loop
			for j in v(i)'Range loop
				v(i)(j) := -1;
			end loop;
		end loop;

	end addPackets;
--------------------------------
	nodes : array (0 .. (n - 1)) of Node;

	procedure addNodes is

		procedure draw(i : Integer; j : Integer) is
			a : Integer := 0;
		begin
			if i < j then
				while a < i loop
					Put("    ");
					a := a + 1; 
				end loop;
				Put(Integer'Image(i));
				while a < j - 1 loop
					Put("----");
					a := a + 1;
				end loop;
				Put("-->" & Integer'Image(j));
				New_Line;
			else
				while a < j loop
					Put("    ");
					a := a + 1; 
				end loop;
				Put(Integer'Image(j) & "<--");
				while a < i - 1 loop
					Put("----");
					a := a + 1;
				end loop;
				Put(Integer'Image(i));
				New_Line;
			end if;
		end draw;

		c : connectionsArray;
		j : Integer;
		l : Integer;
		r : Integer;
		was : Boolean;
		s : servicedPackets;
	begin
		Put_Line("Graph:");
		for i in nodes'Range loop
			j := 0;

			for k in c'Range loop
				c(k) := -1;
			end loop;

			while j < arrayOfRandom(i) loop
				was := false;
				r := n;

				while r >= n or r <= i + 1 loop
					reset(gen);
					r := Integer(random(gen));
				end loop;

				for k in c'Range loop
					if r = c(k) then
						was := true;
					end if;
				end loop;

				if not was then
					c(j) := r;
					draw(i, c(j));
					j := j + 1;
				end if;
			end loop;

			l := j;

			while j < arrayOfReverseRandom(i) + l loop
				was := false;
				r := n;

				while r >= i or r < 0 loop
					reset(gen);
					r := Integer(random(gen)) - 1;
				end loop;

				for k in c'Range loop
					if r = c(k) then
						was := true;
					end if;
				end loop;

				if not was then
					c(j) := r;
					draw(i, c(j));
					j := j + 1;
				end if;
			end loop;

			c(j) := i + 1;
			if i < n - 1 then
				draw(i, c(j));
			end if;

			for k in s'Range loop
				s(k) := -1;
			end loop;

			nodes (i) := (id => i, packet => -1, connections => c, serviced => s, trap => false);
		end loop;
	end addNodes;
--------------------------------
	task sendPackets;
	task body sendPackets is
		r : randRange;
		rf : Float;
	begin
		for i in 0 .. (k - 1) loop
			reset(gen);
			r := random(gen);
			rf :=  Float (r) * 0.6;
			delay Duration (rf);

			while nodes(0).packet /= -1 loop
				delay 0.1;
			end loop;
			nodes(0).packet := packets(i);
			Put_Line ("Sent packet" & Integer'Image(packets(i)));
		end loop;
	end sendPackets;
--------------------------------
	task type receiveAndSend is
		entry Start (j : Integer);
	end receiveAndSend;

	task body receiveAndSend is
		taskJ : Integer;
		r : randRange;
		rf : Float;
		kWas : Boolean := false;
	begin
		accept Start (j : Integer) do
			taskJ := j;
		end Start;
		while true loop
			
			while nodes(taskJ).packet = -1 loop
				delay 0.1;
			end loop;
			Put_Line("Packet" & Integer'Image(nodes(taskJ).packet) & " received by" & Integer'Image(nodes(taskJ).id));

			for k in nodes(taskJ).serviced'Range loop
				if nodes(taskJ).serviced(k) = -1 then
					nodes(taskJ).serviced(k) := nodes(taskJ).packet;
				end if;
				exit when nodes(taskJ).serviced(k) = nodes(taskJ).packet;
			end loop;

			for k in v(nodes(taskJ).packet)'Range loop
				if v(nodes(taskJ).packet)(k) = -1 then
					v(nodes(taskJ).packet)(k) := nodes(taskJ).id;
				end if;
				exit when v(nodes(taskJ).packet)(k) = nodes(taskJ).id;
			end loop;

			if packetLifetime(nodes(taskJ).packet) > h then
				Put_Line("Packet" & Integer'Image(nodes(taskJ).packet) & " reached lifetime in node" & Integer'Image(nodes(taskJ).id));
				nodes(taskJ).packet := -1;
				counter := counter - 1;
				exit when true;
			end if;
			packetLifetime(nodes(taskJ).packet) := packetLifetime(nodes(taskJ).packet) + 1;

			if nodes(taskJ).trap then
				Put_Line("Packet" & Integer'Image(nodes(taskJ).packet) & " has been caught by trap in node" & Integer'Image(nodes(taskJ).id));
				kWas := true;
				goto Continue;
				--exit when true;
			end if;

			reset(gen);
			r := random(gen);
			rf :=  Float (r) * 0.6;
			delay Duration(rf);

			while nodes(taskJ).connections(Integer(r) - 1) = -1 loop
				reset(gen);
				r := random(gen);
			end loop;


			kWas := false;
			while nodes(nodes(taskJ).connections(Integer(r) - 1)).packet /= -1 loop
				delay 0.1;
----------------------------
				if nodes(taskJ).trap then
					Put_Line("Packet" & Integer'Image(nodes(taskJ).packet) & " has been caught by trap in node" & Integer'Image(nodes(taskJ).id));
					kWas := true;
					goto Continue;
				end if;
-----------------------------
			end loop;

			<<Continue>>
			if kWas then
				nodes(taskJ).trap := false;
				nodes(taskJ).packet := -1;
				counter := counter - 1;
				kWas := false;
			else
				nodes(nodes(taskJ).connections(Integer(r) - 1)).packet := nodes(taskJ).packet;
				nodes(taskJ).packet := -1;
			end if;

		end loop;
	end receiveAndSend;

	nTasks : array (0 .. (n - 2)) of receiveAndSend;
--------------------------------

	task lastNode;
	task body lastNode is
		r : randRange;
		rf : Float;
	begin
		delay 1.0;
		for i in 0 .. (k - 1) loop
			while nodes(n - 1).packet = -1 loop
				delay 0.1;
			end loop;
			Put_Line("Packet" & Integer'Image(nodes(n - 1).packet) & " received by" & Integer'Image(nodes(n - 1).id));

			for k in nodes(n - 1).serviced'Range loop
				if nodes(n - 1).serviced(k) = -1 then
					nodes(n - 1).serviced(k) := nodes(n - 1).packet;
				end if;
				exit when nodes(n - 1).serviced(k) = nodes(n - 1).packet;
			end loop;

			for k in v(nodes(n - 1).packet)'Range loop
				if v(nodes(n - 1).packet)(k) = -1 then
					v(nodes(n - 1).packet)(k) := nodes(n - 1).id;
				end if;
				exit when v(nodes(n - 1).packet)(k) = nodes(n - 1).id;
			end loop;

			reset(gen);
			r := random(gen);
			rf :=  Float (r) * 0.6;
			delay Duration(rf);

			lastDone := nodes(n - 1).packet;
			nodes(n - 1).packet := -1;
		end loop;
	end lastNode;
--------------------------------
	task receivePackets;
	task body receivePackets is
		procedure endProgram is
		begin
			New_Line;
			Put_Line("Packets serviced by nodes:");
			for i in nodes'Range loop
				Put(Integer'Image(nodes(i).id) & ": [");
				for j in nodes(i).serviced'Range loop
					if nodes(i).serviced(j) /= -1 then
						Put(Integer'Image(nodes(i).serviced(j)));
					end if;
				end loop;
				Put("]");
				New_Line;
			end loop;

			New_Line;
			Put_Line("Nodes visited by packets:");
			for i in v'Range loop
				Put(Integer'Image(i) & ": [");
				for j in v(i)'Range loop
					if v(i)(j) /= -1 then
						Put(Integer'Image(v(i)(j)));
					end if;
				end loop;
				Put("]");
				New_Line;
			end loop;


			GNAT.OS_Lib.OS_Exit (0);
		end endProgram;

	begin
		while counter > 0 loop
			while lastDone = -1 loop
				if counter <= 0 then
					endProgram;
				end if;
				delay 0.01;
			end loop;
			Put_Line ("Received packet" & Integer'Image(lastDone));
			counter := counter - 1;
			lastDone := -1;
		end loop;

		endProgram;
	end receivePackets;
--------------------------------
	task klusownik is
		entry Start;
	end klusownik;

	task body klusownik is
		r : randRange;
		rf : Float;
	begin
		accept Start do

		delay 1.0;
		while true loop

			reset(gen);
			r := random(gen);
			reset(gen);
			rf :=  Float (random(gen)) * 3.2;
			delay Duration(rf);

			Put_Line("klusownik set trap in node" & Integer'Image(nodes(Integer(r) - 1).id));
			nodes(Integer(r) - 1).trap := true;
		end loop;

		end Start;
	end klusownik;
--------------------------------
	procedure start is
	begin
		checkD;
		checkB;
		generateRandomArray;
		generateReverseRandomArray;
		addPackets;
		addNodes;

		for i in nTasks'Range loop
			nTasks(i).Start (i);
		end loop;

		if kk then
			klusownik.Start;
		end if;
	end start;

begin
	start;
end Main;
